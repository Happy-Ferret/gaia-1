package client

import (
	"io/ioutil"
	"net/http"
)

type Client struct {
	IpAddress      string
	Location       string
	GaiaServerHost string
}

func Start(gaia string) {
	c := NewClient(gaia)
	c.Register()
}

func NewClient(gaia string) *Client {
	ip, location := getPublicIp()
	c := Client{
		IpAddress:      ip,
		Location:       location,
		GaiaServerHost: gaia,
	}

	return &c
}

// Register this client with Gaia server
func (c *Client) Register() {
	// Send post request to gaia server
}

func getPublicIp() (string, string) {
	// @TODO Move this to gaia server
	// request ifconfig.me to find public ip
	resp, err := http.Get("ifconfig.co")
	if err != nil {
		return "", ""
	}

	ip, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return "", ""
	}

	resp, err = http.Get("ifconfig.co/country")
	if err != nil {
		return "", ""
	}

	resp.Body.Close()
	location, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", ""
	}

	return string(ip), string(location)
}
