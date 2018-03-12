package carrefour

import (
	"testing"

	"github.com/Hunsin/bee/mart"
)

const (
	regPage  = "https://online.carrefour.com.tw/[0-9]+"
	regImage = "https://carrefoureccdn.azureedge.net/content/images/thumbs/.+.[jpe?g|png]"
)

func TestSeek(t *testing.T) {
	err := mart.ValidSeek(regPage, regImage, &client{})
	if err != nil {
		t.Error("client.Seek failed:", err)
	}
}
