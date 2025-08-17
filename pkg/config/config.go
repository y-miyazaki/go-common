// Package config provides configuration management utilities.
package config

import (
	"fmt"
	"os"
	"strings"

	"go-common/pkg/infrastructure"
	"go-common/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
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
	config := &Config{
		Logger:      nil,
		SlackClient: nil,
	}

	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	l := &logrus.Logger{
		Out:          os.Stdout,
		Hooks:        make(logrus.LevelHooks),
		Formatter:    nil,
		ReportCaller: false,
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		BufferPool:   nil,
	}

	// formatter
	formatter := strings.ToLower(setting.LoggerFormatter)
	switch formatter {
	case "json":
		l.Formatter = &logrus.JSONFormatter{
			TimestampFormat:   "",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       false,
		}
	case "text":
		l.Formatter = &logrus.TextFormatter{
			ForceColors:               false,
			DisableColors:             false,
			ForceQuote:                false,
			DisableQuote:              false,
			EnvironmentOverrideColors: false,
			DisableTimestamp:          false,
			FullTimestamp:             false,
			TimestampFormat:           "",
			DisableSorting:            false,
			SortingFunc:               nil,
			DisableLevelTruncation:    false,
			PadLevelText:              false,
			QuoteEmptyFields:          false,
			FieldMap:                  nil,
			CallerPrettyfier:          nil,
		}
	default:
		panic("Only json and text can be selected for formatter.")
	}

	// out
	out := strings.ToLower(setting.LoggerOut)
	switch out {
	case "stdout":
		l.Out = os.Stdout
	default:
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
