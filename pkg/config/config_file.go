package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// FileSetting sets configurations.
type FileSetting struct {
	ConfigPath            string
	ConfigFileName        string
	SlackOauthAccessToken string
}

// NewConfigFile to read config
func NewConfigFile(setting *FileSetting) *Config {
	config := &Config{
		Logger:      nil,
		SlackClient: nil,
	}
	// -------------------------------------------------------------
	// get config
	// -------------------------------------------------------------
	viper.AddConfigPath(setting.ConfigPath) // path to look for the config file in
	viper.SetConfigName(setting.ConfigFileName)
	viper.SetConfigType("yaml") // can viper.SetConfigType("YAML")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		panic("viper can't read config.")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// use logger instead of fmt.Println
		log := logrus.New()
		log.Infof("ConfigHandler file changed: %s", e.Name)
	})
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
	formatter := strings.ToLower(viper.GetString("logger.formatter"))
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
	out := strings.ToLower(viper.GetString("logger.out"))
	if out == "stdout" {
		l.Out = os.Stdout
	} else {
		panic("Only stdout can be selected for out.")
	}
	// level
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
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
