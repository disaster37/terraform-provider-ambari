package ambari

import (
	restClient "github.com/disaster37/go-ambari-rest"
)

type Client struct {
	APIURL   string
	Login    string
	Password string
	client   *restClient.AmbariClient
}

func NewClient() *Client {
	return &Client{}
}

// Permit to get rest API client
func (c *Client) Client() *restClient.AmbariClient {

	if c.client == nil {
		if c.APIURL == "" {
			panic("You must provide APIURL")
		}
		if c.Login == "" {
			panic("You must provide Login")
		}
		if c.Password == "" {
			panic("You must provide Password")
		}

		c.client = restClient.New(c.APIURL, c.Login, c.Password)
		c.client.DisableVerifySSL()
	}

	return c.client
}
