package carrefour

import "mart"

const baseURL = "https://online.carrefour.com.tw"

// A client implements the mart.Mart interface.
type client struct{}

func (c *client) ID() string {
	return "carrefour"
}

func (c *client) Name() string {
	return "Carrefour (TW)"
}

func (c *client) Currency() string {
	return "TWD"
}

// init registers a client to package mart.
func init() {
	c := &client{}
	mart.Register(c.ID(), c)
}
