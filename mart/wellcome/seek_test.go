package wellcome

import (
	"testing"

	"github.com/Hunsin/bee/mart"
)

const (
	regPage  = "https://sbd-ec.wellcome.com.tw/product/view/[0-9a-zA-Z]+"
	regImage = "https://sbd-ec.wellcome.com.tw/fileHandler/show/[0-9]+.+"
)

func TestSeek(t *testing.T) {
	err := mart.ValidSeek(regPage, regImage, &client{})
	if err != nil {
		t.Error("client.Seek failed:", err)
	}
}
