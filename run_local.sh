#!/bin/bash

export NEWRELIC_LICENSE=4b1afcecf6cd215fed9d0c3f1604ab9d9e68eb99
export INFLUXDB_HOST=http://docker.ax:8086
export INFLUXDB_DB=notyim
export INFLUXDB_USERNAME=notyim
export INFLUXDB_PASSWORD='notyim'

export RETHINKDB_HOST=docker.ax
export RETHINKDB_PORT=28015
export RETHINKDB_USER=admin
export RETHINKDB_PASS=""
export RETHINKDB_NAME=notyim

export BIND="127.0.0.1:23501"
echo "Run"
./gaia monitor
