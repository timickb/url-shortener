package server

type DbConfig struct {
	DbHost     string `yaml:"db_host"`
	DbUser     string `yaml:"db_user"`
	DbName     string `yaml:"db_name"`
	DbPassword string `yaml:"db_password"`
	DbPort     int    `yaml:"db_port"`
}

type Config struct {
	ServerPort     int      `yaml:"server_port"`
	MaxUrlLength   int      `yaml:"max_url_length"`
	ShorteningSize int      `yaml:"shortening_size"`
	Database       DbConfig `yaml:"database"`
	Origins        []string `yaml:"origins"`
}

func DefaultConfig() *Config {
	return &Config{
		ServerPort:     8080,
		MaxUrlLength:   300,
		ShorteningSize: 10,
	}
}
