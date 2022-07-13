package urlapi

import "github.com/timickb/url-shortener/internal/app/store"

type Config struct {
	ServerAddress string `yaml:"server_address"`
	LogLevel      string `yaml:"log_level"`
	Store         *store.Config
}

func DefaultConfig() *Config {
	return &Config{
		ServerAddress: ":8080",
		LogLevel:      "debug",
		Store:         store.DefaultConfig(),
	}
}
