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

// Setting sets base configurations.
type Setting struct {
	LoggerFormatter       string
	LoggerOut             string
	LoggerLevel           string
	SlackOauthAccessToken string
}

// NewConfig sets base configurations.
func NewConfig(setting *Setting) *Config {
	config := &Config{}

	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	l := &logrus.Logger{}
	// formatter
	formatter := strings.ToLower(setting.LoggerFormatter)
	if formatter == "json" {
		l.Formatter = &logrus.JSONFormatter{}
	} else if formatter == "text" {
		l.Formatter = &logrus.TextFormatter{}
	} else {
		panic("Only json and text can be selected for formatter.")
	}
	// out
	out := strings.ToLower(setting.LoggerOut)
	if out == "stdout" {
		l.Out = os.Stdout
	} else {
		panic("Only stdout can be selected for out.")
	}
	// level
	level, err := logrus.ParseLevel(setting.LoggerLevel)
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	l.Level = level
	config.Logger = logger.NewLogger(l)

	// -------------------------------------------------------------
	// set Slack
	// -------------------------------------------------------------
	if setting.SlackOauthAccessToken != "" {
		config.SlackClient = infrastructure.NewSlack(
			&infrastructure.SlackConfig{
				OauthAccessToken: setting.SlackOauthAccessToken,
			})
	}
	return config
}
