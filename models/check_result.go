package models

import (
	"github.com/notyim/gaia/types"
	"log"
)

type CheckResult types.HTTPCheckResponse

func (c *CheckResult) Save() {
	log.Println("FLush", c, "to InfluxDB")
}
