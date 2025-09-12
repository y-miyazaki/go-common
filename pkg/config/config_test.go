package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	setting := &Setting{
		LoggerFormatter:       "json",
		LoggerOut:             "stdout",
		LoggerLevel:           "info",
		SlackOauthAccessToken: "",
	}

	config := NewConfig(setting)

	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.Nil(t, config.SlackClient)
}

func TestNewConfigWithSlack(t *testing.T) {
	setting := &Setting{
		LoggerFormatter:       "text",
		LoggerOut:             "stdout",
		LoggerLevel:           "debug",
		SlackOauthAccessToken: "test-token",
	}

	config := NewConfig(setting)

	assert.NotNil(t, config)
	assert.NotNil(t, config.Logger)
	assert.NotNil(t, config.SlackClient)
}

func TestNewConfigInvalidFormatter(t *testing.T) {
	setting := &Setting{
		LoggerFormatter:       "invalid",
		LoggerOut:             "stdout",
		LoggerLevel:           "info",
		SlackOauthAccessToken: "",
	}

	assert.Panics(t, func() {
		NewConfig(setting)
	})
}

func TestNewConfigInvalidOut(t *testing.T) {
	setting := &Setting{
		LoggerFormatter:       "json",
		LoggerOut:             "invalid",
		LoggerLevel:           "info",
		SlackOauthAccessToken: "",
	}

	assert.Panics(t, func() {
		NewConfig(setting)
	})
}

func TestNewConfigInvalidLevel(t *testing.T) {
	setting := &Setting{
		LoggerFormatter:       "json",
		LoggerOut:             "stdout",
		LoggerLevel:           "invalid",
		SlackOauthAccessToken: "",
	}

	assert.Panics(t, func() {
		NewConfig(setting)
	})
}
