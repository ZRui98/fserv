package config

import "os"

type Config struct {
	DB_URL string
	ROOT_DIR string
	JWT_KEY []byte
	REGISTRATION_KEY string
}

func LoadConfig() *Config {
	return &Config{
		DB_URL: os.Getenv("DATABASE_URL"),
		REGISTRATION_KEY: os.Getenv("REGISTRATION_KEY"),
		JWT_KEY: []byte(os.Getenv("SECRET_KEY")),
		ROOT_DIR: os.Getenv("FILE_DIR"),
	}
}
