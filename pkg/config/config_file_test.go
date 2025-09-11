package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFile(t *testing.T) {
	// Create a temporary config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `
logger:
  formatter: json
  out: stdout
  level: info
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	setting := &FileSetting{
		ConfigPath:            dir,
		ConfigFileName:        "config",
		SlackOauthAccessToken: "",
	}

	config := NewConfigFile(setting)

	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.Nil(t, config.SlackClient)
}

func TestNewConfigFileWithSlack(t *testing.T) {
	// Create a temporary config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `
logger:
  formatter: text
  out: stdout
  level: debug
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	setting := &FileSetting{
		ConfigPath:            dir,
		ConfigFileName:        "config",
		SlackOauthAccessToken: "test-token",
	}

	config := NewConfigFile(setting)

	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.NotNil(t, config.SlackClient)
}

func TestNewConfigFileInvalidConfig(t *testing.T) {
	setting := &FileSetting{
		ConfigPath:            "/nonexistent",
		ConfigFileName:        "config",
		SlackOauthAccessToken: "",
	}

	assert.Panics(t, func() {
		NewConfigFile(setting)
	})
}
