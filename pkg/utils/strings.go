// utils is a meaningful package name for utility functions
// nolint:revive
package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/rivo/uniseg"
	"golang.org/x/exp/utf8string"
)

// ErrUndefinedTypeToConvert is returned when ConvertToString receives an unsupported type.
var ErrUndefinedTypeToConvert = errors.New("undefined type to convert")

const maxByteMultiplier = 3

// GetStringCount gets the number of characters in a person's view.
func GetStringCount(str string) int {
	return uniseg.GraphemeClusterCount(str)
}

// CheckStringCount checks if the specified string length has been exceeded.
func CheckStringCount(str string, maxLen int) bool {
	gcc := uniseg.GraphemeClusterCount(str)
	if gcc < 1 || maxLen < gcc {
		return false
	}
	// check byte length.
	return len(str) <= maxLen*maxByteMultiplier
}

// SliceUTF8 gets the characters from the beginning to the specified position,
// in UTF-8-based characters.
func SliceUTF8(str string, pos int) string {
	s := utf8string.NewString(str)
	length := GetStringCount(str)
	if pos >= length {
		return s.Slice(0, length)
	}
	return s.Slice(0, pos)
}

// SliceUTF8AddString gets the characters from the beginning to the specified position,
// in UTF-8-based characters.
func SliceUTF8AddString(str string, pos int, addString string) string {
	s := utf8string.NewString(str)
	length := GetStringCount(str)
	if pos >= length {
		return s.Slice(0, length)
	}
	return s.Slice(0, pos) + addString
}

// ConvertToString converts a basic type to a string type.
func ConvertToString(input any) (string, error) {
	var output string
	var err error
	switch v := input.(type) {
	case int:
		output = strconv.Itoa(v)
	case bool:
		output = strconv.FormatBool(v)
	case float32:
		output = fmt.Sprintf("%f", v)
	case float64:
		output = fmt.Sprintf("%f", v)
	case string:
		output = v
	default:
		// Wrap a static error with the concrete type to preserve type info while avoiding dynamic errors.
		err = fmt.Errorf("%w: %T", ErrUndefinedTypeToConvert, v)
	}
	return output, err
}

// GetStringFromReadCloser gets string
func GetStringFromReadCloser(r io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", fmt.Errorf("read from reader: %w", err)
	}
	return buf.String(), nil
}
