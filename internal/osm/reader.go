package osm

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"gpath/internal/pbf"

	"google.golang.org/protobuf/proto"
)

func getBlobHeader(file *os.File, headerLen uint32) (*pbf.BlobHeader, error) {
	blobHeader := &pbf.BlobHeader{}

	blobHeaderBuf := make([]byte, headerLen)
	if _, err := io.ReadFull(file, blobHeaderBuf); err != nil {
		return blobHeader, fmt.Errorf("failed to read BlobHeader bytes: %w", err)
	}

	if err := proto.Unmarshal(blobHeaderBuf, blobHeader); err != nil {
		return blobHeader, fmt.Errorf("failed to unmarshal BlobHeader: %w", err)
	}

	return blobHeader, nil
}

func getRawData(file *os.File, blobHeader *pbf.BlobHeader) ([]byte, error) {
	var rawData []byte

	blobBuf := make([]byte, blobHeader.GetDatasize())
	if _, err := io.ReadFull(file, blobBuf); err != nil {
		return rawData, fmt.Errorf("failed to read Blob bytes: %w", err)
	}

	var blob pbf.Blob
	if err := proto.Unmarshal(blobBuf, &blob); err != nil {
		return rawData, fmt.Errorf("failed to unmarshal Blob: %w", err)
	}

	if zlibData := blob.GetZlibData(); len(zlibData) > 0 {
		r, err := zlib.NewReader(bytes.NewReader(zlibData))
		if err != nil {
			return rawData, fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer r.Close()
		rawData, err = io.ReadAll(r)
		if err != nil {
			return rawData, fmt.Errorf("failed to decompress zlib data: %w", err)
		}
	} else if raw := blob.GetRaw(); len(raw) > 0 {
		rawData = raw
	} else {
		return rawData, fmt.Errorf("unsupported or empty blob data payload")
	}

	return rawData, nil
}

func getBlock(file *os.File) (*pbf.BlobHeader, []byte, error) {
	var headerLen uint32
	if err := binary.Read(file, binary.BigEndian, &headerLen); err != nil {
		if err == io.EOF {
			return &pbf.BlobHeader{}, []byte{}, io.EOF
		}
		return &pbf.BlobHeader{}, []byte{}, fmt.Errorf("failed to read BlobHeader length: %w", err)
	}

	blobHeader, err := getBlobHeader(file, headerLen)
	if err != nil {
		return &pbf.BlobHeader{}, []byte{}, err
	}

	rawData, err := getRawData(file, blobHeader)
	if err != nil {
		return &pbf.BlobHeader{}, []byte{}, err
	}

	return blobHeader, rawData, nil
}
