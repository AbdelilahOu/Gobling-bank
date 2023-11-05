package config

import "os"

type DbConfig struct {
	User     string
	Password string
	Host     string
	Post     string
}

func LoadDbConfig() DbConfig {
	return DbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Post:     os.Getenv("POSTGRES_PORT"),
	}
}
