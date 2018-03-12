// Package rt is an adapter between RT-Mart online store and package mart.
package rt

import "github.com/Hunsin/bee/mart"

const baseURL = "http://www.rt-mart.com.tw/direct/index.php"

// A client implements the mart.Client interface.
type client struct{}

func (c *client) ID() string {
	return "rt"
}

func (c *client) Name() string {
	return "RT-Mart"
}

func (c *client) Currency() string {
	return mart.CurrencyTWD
}

func init() {
	mart.Register(&client{})
}
