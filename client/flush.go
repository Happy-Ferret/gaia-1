package client

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
	Host string
	ch   chan *types.HTTPCheckResponse
	quit chan bool
}

func NewFlusher(to string) *Flusher {
	f := Flusher{
		Host: to,
		ch:   make(chan *types.HTTPCheckResponse, MaxFlusher),
		quit: make(chan bool),
	}

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

// post check result to gaia
func (f *Flusher) Flush(res *types.HTTPCheckResponse) bool {
	_, err := http.PostForm(fmt.Sprintf("%s/check_results/%d", f.Host, res.CheckID),
		url.Values{
			"TotalTime": {fmt.Sprintf("%d", int64(res.TotalTime/time.Millisecond))},
			"TotalSize": {fmt.Sprintf("%d", res.TotalSize)},
		})

	if err != nil {
		log.Println("Fail to flush", res.CheckID, "err", err)
		return false
	}

	return true
}
