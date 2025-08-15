// utils is a meaningful package name for utility functions
// nolint:revive
package utils

import (
	"bytes"
	"fmt"
	"io"
)

// GetBufferFromReadCloser gets buffer
func GetBufferFromReadCloser(r io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("read from reader: %w", err)
	}
	return buf.Bytes(), nil
}
