// Package carrefour is an adapter between Taiwanese Carrefour online
// store and package mart.
package carrefour

import "github.com/Hunsin/bee/mart"

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
	return mart.CurrencyTWD
}

// init registers a client to package mart.
func init() {
	mart.Register(&client{})
}
