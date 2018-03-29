// Package rt is an adapter between RT-Mart online store and package mart.
package rt

import "github.com/Hunsin/bee/mart"

const (
	id      = "rt"
	baseURL = "http://www.rt-mart.com.tw/direct/index.php"
)

// title is the website's title.
var title = "大潤發網路購物中心"

// A client implements the mart.Client interface.
type client struct{}

func (c *client) Info() mart.Info {
	return mart.Info{
		ID:       id,
		Name:     title,
		Currency: mart.CurrencyTWD,
	}
}

func init() {
	mart.Register(&client{})
}
