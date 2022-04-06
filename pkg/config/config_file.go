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

// ConfigFileSetting sets configurations.
type ConfigFileSetting struct {
	ConfigPath            string
	ConfigFileName        string
	SlackOauthAccessToken string
}

// NewConfigFile to read config
func NewConfigFile(c *ConfigFileSetting) *Config {
	config := &Config{}
	// -------------------------------------------------------------
	// get config
	// -------------------------------------------------------------
	viper.AddConfigPath(c.ConfigPath) // path to look for the config file in
	viper.SetConfigName(c.ConfigFileName)
	viper.SetConfigType("yaml") // can viper.SetConfigType("YAML")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		panic("viper can't read config.")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("ConfigHandler file changed:", e.Name)
	})
	// -------------------------------------------------------------
	// set Logger
	// -------------------------------------------------------------
	l := &logrus.Logger{}
	// formatter
	formatter := strings.ToLower(viper.GetString("logger.formatter"))
	if formatter == "json" {
		l.Formatter = &logrus.JSONFormatter{}
	} else if formatter == "text" {
		l.Formatter = &logrus.TextFormatter{}
	} else {
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
	if c.SlackOauthAccessToken != "" {
		config.SlackClient = infrastructure.NewSlack(
			&infrastructure.SlackConfig{
				OauthAccessToken: c.SlackOauthAccessToken,
			})
	}
	return config
}
