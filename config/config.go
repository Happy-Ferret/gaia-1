package config

import (
	"os"
)

// Config struct hold whole configuration
type Config struct {
	DbHost     string
	DbUser     string
	DbPassword string

	InfluxdbHost     string
	InfluxdbUsername string
	InfluxdbPassword string
	InfluxdbDb       string
}

// NewConfig creates a configuration struct with a sane default value
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

	if val := os.Getenv("INFLUXDB_HOST"); val != "" {
		c.InfluxdbHost = val
	} else {
		c.InfluxdbHost = "http://127.0.0.1:8086"
	}

	if val := os.Getenv("INFLUXDB_USERNAME"); val != "" {
		c.InfluxdbUsername = val
	} else {
		c.InfluxdbUsername = "notyim"
	}

	if val := os.Getenv("INFLUXDB_PASSWORD"); val != "" {
		c.InfluxdbPassword = val
	} else {
		c.InfluxdbPassword = "notyim"
	}

	if val := os.Getenv("INFLUXDB_DB"); val != "" {
		c.InfluxdbDb = val
	} else {
		c.InfluxdbDb = "notyim"
	}

	return c
}
