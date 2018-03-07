package carrefour

import (
	"mart"
	"testing"
)

const (
	regPage  = "https://online.carrefour.com.tw/[0-9]+"
	regImage = "https://carrefoureccdn.azureedge.net/content/images/thumbs/.+.jpeg"
)

func TestSeek(t *testing.T) {
	err := mart.ValidSeek(regPage, regImage, &client{})
	if err != nil {
		t.Error("client.Seek failed:", err)
	}
}