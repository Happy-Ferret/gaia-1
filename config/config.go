package config

import (
	"os"
)

type Config struct {
	DbHost     string
	DbUser     string
	DbPassword string
}

func NewConfig() *Config {
	c := &Config{}

	if val := os.Getenv("DbHost"); val != "" {
		c.DbHost = val
	} else {
		c.DbHost = "127.0.0.1"
	}

	if val := os.Getenv("DbUser"); val != "" {
		c.DbUser = val
	} else {
		c.DbUser = "value"
	}

	return c
}
