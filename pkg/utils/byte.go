package utils

import (
	"bytes"
	"io"
)

// GetBufferFromReadCloser gets buffer
func GetBufferFromReadCloser(r io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
