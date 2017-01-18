package server

import (
	"fmt"
	"github.com/notyim/gaia/types"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	MaxFlusher = 100
)

type Flusher struct {
	ch   chan *types.HTTPCheckResponse
	quit chan bool
}

func NewFlusher() *Flusher {
	f := Flusher{
		ch:   make(chan *types.HTTPCheckResponse, MaxFlusher),
		quit: make(chan bool),
	}

	f.Start()
	return &f
}

// register writter that listens on result channel and flush to gaia
func (f *Flusher) Start() {
	for i := 0; i < MaxFlusher; i++ {
		go func() {
			for {
				select {
				case checkResponse := <-f.ch:
					f.Flush(checkResponse)
				case <-f.quit:
					return
				}
			}
		}()
	}
}

// Write the result to gaia
// This actually queues the check into channel to be flush later
func (f *Flusher) Write(res *types.HTTPCheckResponse) {
	f.ch <- res
}

// post check result to gaia
func (f *Flusher) Flush(res *types.HTTPCheckResponse) bool {
	endpoint := fmt.Sprintf("check_results/%s", res.CheckID)
	log.Println("Flush check result", res, "to", endpoint)

	_, err := http.PostForm(endpoint,
		url.Values{
			"Error":        {fmt.Sprintf("%t", res.Error)},
			"ErrorMessage": {fmt.Sprintf("%s", res.ErrorMessage)},
			"TotalTime":    {fmt.Sprintf("%d", int64(res.TotalTime/time.Millisecond))},
			"TotalSize":    {fmt.Sprintf("%d", res.TotalSize)},
		})

	if err != nil {
		log.Println("Fail to flush", res.CheckID, "err", err)
		return false
	}

	return true
}
