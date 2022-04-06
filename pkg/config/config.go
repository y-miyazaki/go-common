package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// Config sets base configurations.
type Config struct {
	Logger      *logger.Logger
	SlackClient *slack.Client
}

// ConfigSetting sets base configurations.
type ConfigSetting struct {
	LoggerFormatter       string
	LoggerOut             string
	LoggerLevel           string
	SlackOauthAccessToken string
}

// NewConfig sets base configurations.
func NewConfig(c *ConfigSetting) *Config {
	config := &Config{}

	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	l := &logrus.Logger{}
	// formatter
	formatter := strings.ToLower(c.LoggerFormatter)
	if formatter == "json" {
		l.Formatter = &logrus.JSONFormatter{}
	} else if formatter == "text" {
		l.Formatter = &logrus.TextFormatter{}
	} else {
		panic("Only json and text can be selected for formatter.")
	}
	// out
	out := strings.ToLower(c.LoggerOut)
	if out == "stdout" {
		l.Out = os.Stdout
	} else {
		panic("Only stdout can be selected for out.")
	}
	// level
	level, err := logrus.ParseLevel(c.LoggerLevel)
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	l.Level = level
	config.Logger = logger.NewLogger(l)

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
