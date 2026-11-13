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

func parseBlock(file *os.File) error {
	// 1. Read the 4-byte big-endian length of the BlobHeader
	var headerLen uint32
	if err := binary.Read(file, binary.BigEndian, &headerLen); err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read BlobHeader length: %w", err)
	}

	// 2. Read and unmarshal the BlobHeader
	blobHeaderBuf := make([]byte, headerLen)
	if _, err := io.ReadFull(file, blobHeaderBuf); err != nil {
		return fmt.Errorf("failed to read BlobHeader bytes: %w", err)
	}

	var blobHeader pbf.BlobHeader
	if err := proto.Unmarshal(blobHeaderBuf, &blobHeader); err != nil {
		return fmt.Errorf("failed to unmarshal BlobHeader: %w", err)
	}

	// 3. Read and unmarshal the Blob payload
	blobBuf := make([]byte, blobHeader.GetDatasize())
	if _, err := io.ReadFull(file, blobBuf); err != nil {
		return fmt.Errorf("failed to read Blob bytes: %w", err)
	}

	var blob pbf.Blob
	if err := proto.Unmarshal(blobBuf, &blob); err != nil {
		return fmt.Errorf("failed to unmarshal Blob: %w", err)
	}

	// 4. Extract data (handles zlib-compressed or raw blobs)
	var rawData []byte
	if zlibData := blob.GetZlibData(); len(zlibData) > 0 {
		r, err := zlib.NewReader(bytes.NewReader(zlibData))
		if err != nil {
			return fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer r.Close()
		rawData, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to decompress zlib data: %w", err)
		}
	} else if raw := blob.GetRaw(); len(raw) > 0 {
		rawData = raw
	} else {
		return fmt.Errorf("unsupported or empty blob data payload")
	}

	// 5. If it's an OSMHeader, parse it into the HeaderBlock structure
	if blobHeader.GetType() == "OSMHeader" {
		var headerBlock pbf.HeaderBlock
		if err := proto.Unmarshal(rawData, &headerBlock); err != nil {
			return fmt.Errorf("failed to unmarshal HeaderBlock: %w", err)
		}
		fmt.Printf("OSM HeaderBlock parsed successfully!\nRequired Features: %v\n", headerBlock.GetRequiredFeatures())
	} else {
		fmt.Printf("First block type is %s (expected OSMHeader)\n", blobHeader.GetType())
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
