package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

func Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return parseHeader(file)
}

func parseHeader(r io.Reader) error {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("failed to read header length: %w", err)
	}

	headerLen := binary.BigEndian.Uint32(lenBuf)
	fmt.Printf("header_len: %d\n", headerLen)

	headerData := make([]byte, headerLen)
	if _, err := io.ReadFull(r, headerData); err != nil {
		return fmt.Errorf("unexpected EOF on BlobHeader data: %w", err)
	}

	return nil
}

