package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

type Client struct {
	IpAddress      string
	Location       string
	GaiaServerHost string
}

// Initliaze gaia client
// Register client with the server, create scanner to run check
// and an HTTP server to interact with Gaia server
func Start(gaia string) {
	log.Println("Initalize client with gaia at", gaia)
	c := NewClient(gaia)
	log.Println("Register client")
	c.Register()

	log.Println("Create scanner")
	s := NewScanner(c.GaiaServerHost)
	s.Start()

	log.Println("Create http listener")
	CreateHTTPServer(s)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	userSignal := <-sigChan
	log.Println("Got signal", userSignal, "from end-user")
	log.Println("Attempt to quit")
	os.Exit(0)
}

func NewClient(gaia string) *Client {
	ip, location := getGeoIP()
	log.Println("Found", ip, location)

	// TODO abstract this
	if os.Getenv("ENV") == "LOCAL" {
		ip = "127.0.0.1"
	}

	c := Client{
		IpAddress:      ip,
		Location:       location,
		GaiaServerHost: gaia,
	}

	return &c
}

// Register this client with Gaia server
func (c *Client) Register() {
	_, err := http.PostForm(fmt.Sprintf("%s/client/register", c.GaiaServerHost),
		url.Values{"ip": {c.IpAddress}, "location": {c.Location}})
	log.Println("Register myself as", c.IpAddress, "at", c.Location)
	if err != nil {
		log.Fatal("Fail to register client", err)
	}
}

func getGeoIP() (string, string) {
	// @TODO Move this to gaia server
	// request ifconfig.me to find public ip
	resp, err := http.Get("https://ifconfig.co")
	if err != nil {
		log.Println("Error fetch geoip", err)
		return "", ""
	}

	ip, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Println("Error fetch read ip", err)
		return "", ""
	}

	resp, err = http.Get("https://ifconfig.co/country")
	if err != nil {
		log.Println("Error fetch geo location", err)
		return "", ""
	}

	location, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("Error read geo location", err)
		return "", ""
	}

	return string(ip), string(location)
}
