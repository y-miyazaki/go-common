package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// LambdaConfig sets lambda configurations.
type LambdaConfig struct {
	Logger      *infrastructure.Logger
	SlackClient *slack.Client
}

// ConfigSetting sets lambda configurations.
type ConfigSetting struct {
	LoggerFormatter       string
	LoggerOut             string
	LoggerLevel           string
	SlackOauthAccessToken string
}

// NewConfig sets lambda configurations.
func NewConfig(c *ConfigSetting) *LambdaConfig {
	config := &LambdaConfig{}

	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	logger := &logrus.Logger{}
	// formatter
	formatter := strings.ToLower(c.LoggerFormatter)
	if formatter == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	} else if formatter == "text" {
		logger.Formatter = &logrus.TextFormatter{}
	} else {
		panic("Only json and text can be selected for formatter.")
	}
	// out
	out := strings.ToLower(c.LoggerOut)
	if out == "stdout" {
		logger.Out = os.Stdout
	} else {
		panic("Only stdout can be selected for out.")
	}
	// level
	level, err := logrus.ParseLevel(c.LoggerLevel)
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	logger.Level = level
	config.Logger = infrastructure.NewLogger(logger)

	// -------------------------------------------------------------
	// set Slack
	// -------------------------------------------------------------
	if c.SlackOauthAccessToken != "" {
		config.SlackClient = infrastructure.NewSlack(
			&infrastructure.SlackConfig{
				OauthAccessToken: c.SlackOauthAccessToken,
			})
	}
	return config
}
