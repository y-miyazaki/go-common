package utils

import (
	"fmt"
	"strconv"

	"github.com/rivo/uniseg"
)

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
	return len(str) <= maxLen*3
}

// ConvertToString converts a basic type to a string type.
func ConvertToString(input interface{}) (string, error) {
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
		err = fmt.Errorf("Undefined type to convert %T", v)
	}
	return output, err
}
