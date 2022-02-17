package infrastructure

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// LambdaConfig sets lambda configurations.
type LambdaConfig struct {
	Logger      *Logger
	SlackClient *SlackClient
}

// LambdaConfigSetting sets lambda configurations.
type LambdaConfigSetting struct {
	LoggerFormatter       string
	LoggerOut             string
	LoggerLevel           string
	SlackOauthAccessToken string
	SlackChannelID        string
}

// NewLambdaConfig sets lambda configurations.
func NewLambdaConfig(lcs *LambdaConfigSetting) *LambdaConfig {
	config := &LambdaConfig{}
	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	logger := &logrus.Logger{}
	// formatter
	formatter := strings.ToLower(lcs.LoggerFormatter)
	if formatter == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	} else if formatter == "text" {
		logger.Formatter = &logrus.TextFormatter{}
	} else {
		panic("Only json and text can be selected for formatter.")
	}
	// out
	out := strings.ToLower(lcs.LoggerOut)
	if out == "stdout" {
		logger.Out = os.Stdout
	} else {
		panic("Only stdout can be selected for out.")
	}
	// level
	level, err := logrus.ParseLevel(lcs.LoggerLevel)
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	logger.Level = level
	config.Logger = NewLogger(logger)

	// -------------------------------------------------------------
	// set Slack
	// -------------------------------------------------------------
	if lcs.SlackOauthAccessToken != "" && lcs.SlackChannelID != "" {
		config.SlackClient = NewSlack(lcs.SlackOauthAccessToken, lcs.SlackChannelID)
	}
	return config
}
