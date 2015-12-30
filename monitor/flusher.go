package monitor

import (
	"log"
)

type Flusher struct {
	DataChan chan StatusResult
	Size     int
}

func NewFlusher() *Flusher {
	f := &Flusher{}
	f.Size = 1000
	f.DataChan = make(chan StatusResult, f.Size)
	return f
}

func (f *Flusher) Start() {
	for {
		r := <-f.DataChan
		log.Printf("Will flush result %v", r)
	}
}
