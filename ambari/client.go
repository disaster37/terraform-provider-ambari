package ambari

import (
	"github.com/dghubble/sling"
	"net/http"
)

type Client struct {
	APIURL   string
	Login    string
	Password string
	client   *Sling
}

func NewClient() *Client {
	return &Client{}
}

// Permit to get rest API client
func (c *Client) Client() *Sling {
	if c.client == nil {
		if APIURL == "" {
			panic("You must provide APIURL")
		}
		if Login == "" {
			panic("You must provide Login")
		}
		if Password == "" {
			panic("You must provide Password")
		}

		c.client = sling.New().Base(c.APIURL).Set("X-Requested-By", "ambari").SetBasicAuth(c.Login, c.Password)
	}

	return c.client
}
