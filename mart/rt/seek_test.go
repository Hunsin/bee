package rt

import (
	"testing"

	"github.com/Hunsin/bee/mart/marttest"
)

const (
	regPage  = "http://www.rt-mart.com.tw/direct/index.php.action=product_detail&prod_no=[A-Z0-9]+"
	regImage = "http://www.rt-mart.com.tw/website/uploads_product/website_[0-9]+/.+.[jpe?g|png]"
)

func TestSeek(t *testing.T) {
	err := marttest.ValidSeek(regPage, regImage, &client{})
	if err != nil {
		t.Error("client.Seek failed:", err)
	}
}
