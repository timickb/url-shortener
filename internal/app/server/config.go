package server

type DbConfig struct {
	DbHost     string `yaml:"db_host"`
	DbUser     string `yaml:"db_user"`
	DbName     string `yaml:"db_name"`
	DbPassword string `yaml:"db_password"`
	DbPort     int    `yaml:"db_port"`
}

type Config struct {
	ServerAddress string   `yaml:"server_address"`
	MaxUrlLength  int      `yaml:"max_url_length"`
	Database      DbConfig `yaml:"database"`
	StoreImpl     string
}

func DefaultConfig() *Config {
	return &Config{
		ServerAddress: ":8080",
		StoreImpl:     "db",
		MaxUrlLength:  120,
	}
}
