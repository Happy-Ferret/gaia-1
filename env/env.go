package env

import (
	"github.com/influxdb/influxdb/client/v2"
	"github.com/notyim/gaia/config"
)

type Env struct {
	Config *config.Config
	Influx client.Client
}

var (
	app *Env
)

func NewEnv() *Env {
	config := config.NewConfig()
	app = &Env{
		Config: config,
	}

	app.initInflux()
	app.initRethink()
	return app
}

func GetEnv() *Env {
	return app
}

func (e *Env) initInflux() error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     e.Config.InfluxdbHost,
		Username: e.Config.InfluxdbUsername,
		Password: e.Config.InfluxdbPassword,
	})
	if err != nil {
		return err
	}
	e.Influx = c
	return nil
}

func (e *Env) initRethink() error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     e.Config.InfluxdbHost,
		Username: e.Config.InfluxdbUsername,
		Password: e.Config.InfluxdbPassword,
	})
	if err != nil {
		return err
	}
	e.Influx = c
	return nil
}
