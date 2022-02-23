package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// LambdaConfig sets lambda configurations.
type LambdaConfig struct {
	Logger      *infrastructure.Logger
	SlackClient *slack.Client
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
func NewLambdaConfig(c *LambdaConfigSetting) *LambdaConfig {
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
			&infrastructure.SlackConfigSetting{
				OauthAccessToken: c.SlackOauthAccessToken,
			})
	}
	return config
}
