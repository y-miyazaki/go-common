package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/y-miyazaki/go-common/pkg/errors"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// NewConfig to read config
func NewConfig(configPath string, fileName string) error {
	viper.AddConfigPath(configPath) // path to look for the config file in
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml") // can viper.SetConfigType("YAML")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return errors.Wrap("viper can't read config.", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("ConfigHandler file changed:", e.Name)
	})

	return nil
}

// SetConfig set value to config file.
func SetConfig(key string, value interface{}) {
	viper.Set(key, value)
}

// GetConfigString get string from config file.
func GetConfigString(key string) string {
	return viper.GetString(key)
}

// GetConfigInt get int from config file.
func GetConfigInt(key string) int {
	return viper.GetInt(key)
}

// GetConfigInt64 get int64 from config file.
func GetConfigInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetConfigBool get bool from config file.
func GetConfigBool(key string) bool {
	return viper.GetBool(key)
}

// GetConfigStringMap get bool from config file.
func GetConfigStringMap(key string) interface{} {
	return viper.GetStringMap(key)
}

// GetConfigByte get []byte from config file.
func GetConfigByte(key string) []byte {
	return []byte(viper.GetString(key))
}

// GetConfigEnv from env file.
func GetConfigEnv(key string) string {
	return os.Getenv(key)
}

// GetConfigTime get time from config file.
func GetConfigDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
