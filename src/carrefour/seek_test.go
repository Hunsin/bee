package carrefour

import (
	"mart"
	"testing"
)

func TestSeek(t *testing.T) {
	c := &client{}
	ps, pages, err := c.Seek("抽取衛生紙", 1, mart.ByPrice)
	if err != nil {
		t.Error("client.Seek failed:", err)
	}

	t.Log(pages)
	for i := range ps {
		t.Log(ps[i])
	}
}
