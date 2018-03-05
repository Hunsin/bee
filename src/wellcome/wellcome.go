package wellcome

import "mart"

const (
	pkgName = "wellcome"
	baseURL = "https://sbd-ec.wellcome.com.tw"
)

// A client implements the mart.Mart interface.
type client struct{}

// init registers a client to package mart.
func init() {
	mart.Register(pkgName, &client{})
}
