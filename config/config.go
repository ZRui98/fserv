package config

import (
	"os"
	"github.com/golang/glog"
)

type Config struct {
	DB_URL string
	ROOT_DIR string
	JWT_KEY []byte
	REGISTRATION_KEY string
}

func LoadConfig() *Config {
	dbUrl :=os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		glog.Fatal("DB_URL not specified in environment!")
		os.Exit(1)
	}
	rootDir :=os.Getenv("ROOT_DIR")
	if len(rootDir) == 0 {
		glog.Fatal("ROOT_DIR not specified in environment!")
		os.Exit(1)
	}
	registrationKey :=os.Getenv("REGISTRATION_KEY")
	if len(registrationKey) == 0 {
		glog.Fatal("REGISTRATION_KEY not specified in environment!")
		os.Exit(1)
	}
	secretKey :=os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		glog.Fatal("SECRET_KEY not specified in environment!")
		os.Exit(1)
	}
	return &Config{
		DB_URL: dbUrl,
		ROOT_DIR: rootDir,
		REGISTRATION_KEY: registrationKey,
		JWT_KEY: []byte(secretKey),
	}
}
