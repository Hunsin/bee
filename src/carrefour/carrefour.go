package carrefour

import "mart"

const (
	pkgName = "carrefour"
	baseURL = "https://online.carrefour.com.tw"
)

// A client implements the mart.Mart interface.
type client struct{}

// init registers a client to package mart.
func init() {
	mart.Register(pkgName, &client{})
}
