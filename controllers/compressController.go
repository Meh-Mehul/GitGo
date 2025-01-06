package controllers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func CompressData(data []byte) ([]byte, error){
	var compressData bytes.Buffer
	writer := gzip.NewWriter(&compressData)
	defer writer.Close()
	_, err := writer.Write(data)
	if(err != nil){
		return nil,err;
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	return compressData.Bytes(), nil;
}
func DecompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer reader.Close()
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %w", err)
	}
	return decompressedData, nil
}