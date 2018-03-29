// Package carrefour is an adapter between Taiwanese Carrefour online
// store and package mart.
package carrefour

import "github.com/Hunsin/bee/mart"

const (
	id      = "carrefour"
	baseURL = "https://online.carrefour.com.tw"
)

// title is the website's title.
var title = "家樂福線上購物網"

// A client implements the mart.Mart interface.
type client struct{}

func (c *client) Info() mart.Info {
	return mart.Info{
		ID:       id,
		Name:     title,
		Currency: mart.CurrencyTWD,
	}
}

// init registers a client to package mart.
func init() {
	mart.Register(&client{})
}
