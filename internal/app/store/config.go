package store

type Config struct {
	ConnectionString string `yaml:"connection_string"`
}

func DefaultConfig() *Config {
	return &Config{}
}
