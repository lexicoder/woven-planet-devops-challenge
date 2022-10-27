package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	ListenAddress string
	WriteTimeout  time.Duration
	ReadTimeout   time.Duration
	IdleTimeout   time.Duration
	StoragePath   string
}

func NewConfig() *Config {
	var err error
	var readTimeout time.Duration
	var writeTimeout time.Duration

	readTimeout, err = getDurationFromEnv("READ_TIMEOUT", "15s")
	if err != nil {
		panic(err)
	}
	writeTimeout, err = getDurationFromEnv("WRITE_TIMEOUT", "15s")
	if err != nil {
		panic(err)
	}

	cfg := Config{
		ListenAddress: fmt.Sprintf("0.0.0.0:%s", getEnvOrDefault("SERVER_PORT", "5000")),
		ReadTimeout:   readTimeout,
		WriteTimeout:  writeTimeout,
		StoragePath:   getEnvOrDefault("STORAGE_PATH", "/storage/data"),
	}
	return &cfg
}

func getEnvOrDefault(name, value string) string {
	v := os.Getenv(name)
	if v == "" {
		v = value
	}
	return v
}

func getDurationFromEnv(key, value string) (time.Duration, error) {
	v := getEnvOrDefault(key, value)
	return time.ParseDuration(v)
}
