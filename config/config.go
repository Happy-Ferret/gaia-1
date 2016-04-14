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

	RethinkDBHost string
	RethinkDBUser string
	RethinkDBPass string
	RethinkDBName string
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

	if val := os.Getenv("RETHINKDB_HOST"); val != "" {
		c.RethinkDBHost = val
	} else {
		c.InfluxdbDb = "127.0.0.1"
	}
	if val := os.Getenv("RETHINKDB_USER"); val != "" {
		c.RethinkDBUser = val
	} else {
		c.InfluxdbDb = "127.0.0.1"
	}
	if val := os.Getenv("RETHINKDB_PASS"); val != "" {
		c.RethinkDBPass = val
	} else {
		c.InfluxdbDb = "127.0.0.1"
	}
	if val := os.Getenv("RETHINKDB_NAME"); val != "" {
		c.RethinkDBName = val
	} else {
		c.InfluxdbDb = "127.0.0.1"
	}

	return c
}
