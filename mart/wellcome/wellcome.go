package wellcome

import "github.com/Hunsin/bee/mart"

const (
	id      = "wellcome"
	baseURL = "https://sbd-ec.wellcome.com.tw"
)

// title is the website's title.
var title = "頂好新鮮GO"

// A client implements the mart.Mart interface.
type client struct{}

func (c *client) Info() mart.Info {
	return mart.Info{
		ID:       id,
		Name:     title,
		Currency: mart.CurrencyTWD,
	}
}

func (c *client) ID() string {
	return id
}

// init registers a client to package mart.
func init() {
	mart.Register(&client{})
}
