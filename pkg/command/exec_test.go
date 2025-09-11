package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinedOutput(t *testing.T) {
	// Test with a simple command that should work on most systems
	out, err := CombinedOutput("echo", false, "hello")
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(out), "hello")
}

func TestCombinedOutput_WithOutput(t *testing.T) {
	// Test with output enabled (will print to log)
	out, err := CombinedOutput("echo", true, "test")
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(out), "test")
}

func TestCombinedOutput_InvalidCommand(t *testing.T) {
	// Test with invalid command
	_, err := CombinedOutput("nonexistentcommand", false)
	assert.Error(t, err)
}

func TestCombinedOutputStr(t *testing.T) {
	// Test with a simple shell command
	out, err := CombinedOutputStr("echo hello", false)
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(out), "hello")
}

func TestCombinedOutputStr_WithOutput(t *testing.T) {
	// Test with output enabled
	out, err := CombinedOutputStr("echo test", true)
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(out), "test")
}

func TestCombinedOutputStr_InvalidCommand(t *testing.T) {
	// Test with invalid shell command
	_, err := CombinedOutputStr("nonexistentcommand", false)
	assert.Error(t, err)
}
