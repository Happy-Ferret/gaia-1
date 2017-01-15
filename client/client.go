package client

import (
	"net/http"
)

type Client struct {
	IpAddress string
	Location  string
}

func Start() {
	c := NewClient()
	c.Register()
}

func NewClient() *Client {
	ip, location := getPublicIp()
	c := Client{
		IpAddress: ip,
		Location:  location,
	}

	return &c
}

// Register this client with Gaia server
func (c *Client) Register() {
	// Send post request to gaia server
}

func getPublicIp() (ip, location string) {
	// @TODO Move this to gaia server
	// request ifconfig.me to find public ip
	return "8.8.8.8", "France"
}
