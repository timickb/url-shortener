package server

type Config struct {
	ServerAddress    string `yaml:"server_address"`
	ConnectionString string `yaml:"connection_string"`
	StoreImpl        string
	MaxUrlLength     int
}

func DefaultConfig() *Config {
	return &Config{
		ServerAddress: ":8080",
		StoreImpl:     "local",
		MaxUrlLength:  120,
	}
}
