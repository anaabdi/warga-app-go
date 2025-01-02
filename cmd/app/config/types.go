package config

import (
	"os"
	"time"
)

type Config struct {
	AppName           string
	ServerHost        string
	ServerPort        string
	ReadHeaderTimeout time.Duration
}

func NewConfig() *Config {
	appName := os.Getenv("APP_NAME")
	serverPort := os.Getenv("SERVER_PORT")
	serverHost := os.Getenv("SERVER_HOST")

	readHeaderTimeout, err := time.ParseDuration(os.Getenv("READ_HEADER_TIMEOUT"))
	if err != nil {
		readHeaderTimeout = 5 * time.Second
	}

	return &Config{
		AppName:           appName,
		ServerPort:        serverPort,
		ServerHost:        serverHost,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}
