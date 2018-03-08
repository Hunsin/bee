package wellcome

import "mart"

const baseURL = "https://sbd-ec.wellcome.com.tw"

// A client implements the mart.Mart interface.
type client struct{}

func (c *client) ID() string {
	return "wellcome"
}

func (c *client) Name() string {
	return "Wellcome (TW)"
}

func (c *client) Currency() string {
	return mart.CurrencyTWD
}

// init registers a client to package mart.
func init() {
	mart.Register(&client{})
}
