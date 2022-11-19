package musicbrainz

import (
	"fmt"

	"github.com/axatol/anywhere/config"
)

const baseURL = "https://musicbrainz.org/ws/2"

var client *Client

type Pagination struct {
	Created string `json:"created"`
	Count   int    `json:"count"`
	Offset  int    `json:"offset"`
}

type Client struct {
	AppName    string
	AppVersion string
	Contact    string
}

func NewClient() (*Client, error) {
	if client != nil {
		return client, nil
	}

	client := &Client{
		AppName:    "anywhere",
		AppVersion: "0.0.1",
		Contact:    config.Values.Server.AllowOrigins[0],
	}

	return client, nil
}

func (c *Client) UserAgent() string {
	return fmt.Sprintf("%s/%s ( %s )", c.AppName, c.AppVersion, c.Contact)
}
