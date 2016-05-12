package env

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/notyim/gaia/config"
	"log"
)

type Env struct {
	Config  *config.Config
	Influx  client.Client
	Rethink *r.Session
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
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:%s", e.Config.RethinkDBHost, e.Config.RethinkDBPort),
		Database: e.Config.RethinkDBName,
		MaxIdle:  10,
		MaxOpen:  10,
		Username: e.Config.RethinkDBUser,
		Password: e.Config.RethinkDBPass,
	})
	if err != nil {
		log.Fatalf("Cannot connect to RethinkDB %s", err)
		return err
	}

	session.SetMaxOpenConns(10)
	e.Rethink = session
	return nil
}
