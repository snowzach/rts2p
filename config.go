package main

import (
	"strings"

	config "github.com/spf13/viper"
)

func init() {
	// Sets up the config file, environment etc
	config.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	config.AutomaticEnv()                                   // Automatically use environment variables where available
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	// Logger Defaults
	config.SetDefault("server.port", 5554)
	config.SetDefault("server.http_port", "")
	config.SetDefault("server.username", "")
	config.SetDefault("server.password", "")
	config.SetDefault("server.max_out_packet_size", 2000000)

	config.SetDefault("streams", []Stream{})
}
