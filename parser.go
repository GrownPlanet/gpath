package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"gpath/pbf"

	"google.golang.org/protobuf/proto"
)

func getBlobHeader(file *os.File, headerLen uint32) (pbf.BlobHeader, error) {
	var blobHeader pbf.BlobHeader

	blobHeaderBuf := make([]byte, headerLen)
	if _, err := io.ReadFull(file, blobHeaderBuf); err != nil {
		return blobHeader, fmt.Errorf("failed to read BlobHeader bytes: %w", err)
	}

	if err := proto.Unmarshal(blobHeaderBuf, &blobHeader); err != nil {
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


func parseBlock(file *os.File) error {
	var headerLen uint32
	if err := binary.Read(file, binary.BigEndian, &headerLen); err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read BlobHeader length: %w", err)
	}

	blobHeader, err := getBlobHeader(file, headerLen)
	if err != nil {
		return err
	}

	rawData, err := getRawData(file, &blobHeader)
	if err != nil {
		return err
	}

	headerType := blobHeader.GetType();
	fmt.Println(headerType)
	if headerType == "OSMHeader" {
		var headerBlock pbf.HeaderBlock
		if err := proto.Unmarshal(rawData, &headerBlock); err != nil {
			return fmt.Errorf("failed to unmarshal HeaderBlock: %w", err)
		}
		fmt.Printf("OSM HeaderBlock parsed successfully!\nRequired Features: %v\n", headerBlock.GetRequiredFeatures())
	} else if headerType == "OSMData" {
		var primitiveBlock pbf.PrimitiveBlock
		if err := proto.Unmarshal(rawData, &primitiveBlock); err != nil {
			return fmt.Errorf("failed to unmarshal PrimitiveBlockBlock: %w", err)
		}
		fmt.Println(primitiveBlock.String())
	} else {
		return fmt.Errorf("unknown blobHeader type: %s", headerType)
	}
	fmt.Println(blobHeader.GetIndexdata(), blobHeader.GetDatasize())

	return nil
}

func Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	for {
		err = parseBlock(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}
