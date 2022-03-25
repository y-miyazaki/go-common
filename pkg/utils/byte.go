package utils

import (
	"bytes"
	"io"
)

// GetBufferFromReadCloser gets buffer
func GetBufferFromReadCloser(r io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes()
}
