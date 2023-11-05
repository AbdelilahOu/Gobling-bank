package config

type Config struct {
	DB DbConfig
}

func NewConfig() *Config {
	return &Config{
		DB: LoadDbConfig(),
	}
}
