package config

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

// Config struct used for environment variables
type Config struct {
	DB_URL           string
	ROOT_DIR         string
	JWT_KEY          []byte
	REGISTRATION_KEY string
	SECURE_COOKIE    bool
}

// LoadConfig generates config struct from environment variables
func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if len(dbURL) == 0 {
		glog.Fatal("DB_URL not specified in environment!")
		os.Exit(1)
	}
	rootDir := os.Getenv("ROOT_DIR")
	if len(rootDir) == 0 {
		glog.Fatal("ROOT_DIR not specified in environment!")
		os.Exit(1)
	}
	registrationKey := os.Getenv("REGISTRATION_KEY")
	if len(registrationKey) == 0 {
		glog.Fatal("REGISTRATION_KEY not specified in environment!")
		os.Exit(1)
	}
	secretKey := os.Getenv("SECRET_KEY")
	if len(secretKey) == 0 {
		glog.Fatal("SECRET_KEY not specified in environment!")
		os.Exit(1)
	}
	val := os.Getenv("SECURE_COOKIE")
	secureCookie := true
	if len(val) != 0 {
		var err error
		secureCookie, err = strconv.ParseBool(val)
		if err != nil {
			glog.Fatal("SECURE_COOKIE was invalid: Must be one of \"true\" or \"false\"")
		}
	}
	return &Config{
		DB_URL:           dbURL,
		ROOT_DIR:         rootDir,
		REGISTRATION_KEY: registrationKey,
		JWT_KEY:          []byte(secretKey),
		SECURE_COOKIE:    secureCookie,
	}
}
