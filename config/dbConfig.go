package config

import (
	"os"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Post     int
}
type Config struct {
	Db DbConfig
}

func LoadDbConfig() DbConfig {
	var dbCfg DbConfig

	if password, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists {
		dbCfg.Password = password
	} else {
		dbCfg.Password = "mysecretpassword"
	}

	if user, exists := os.LookupEnv("POSTGRES_USER"); exists {
		dbCfg.User = user
	} else {
		dbCfg.User = "root"
	}

	if host, exists := os.LookupEnv("POSTGRES_HOST"); exists {
		dbCfg.Host = host
	} else {
		dbCfg.User = "localhost"
	}

	dbCfg.Post = 5432

	return dbCfg
}

func LoadConfig() *Config {
	dbCfg := LoadDbConfig()
	return &Config{
		Db: dbCfg,
	}
}
