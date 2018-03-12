package wellcome

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mart"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
)

const (
	searchURL = baseURL + "/product/listByKeyword"

	// max numer of items per search
	searchSize = 100

	// product page, image & name
	tmplInfo = "<a href=\"(.*)\">\n *<img src=\"(.*)\" alt=\"(.*)\" class=\"item-image\">"

	// product price
	tmplPrice = `<span class="item-price ">([0-9]*)</span>`

	// number of products
	tmplNum = `<li class="active">關鍵字: .* \(([0-9]*)\)</li>`
)

var (
	regInfo  = regexp.MustCompile(tmplInfo)
	regPrice = regexp.MustCompile(tmplPrice)
	regNum   = regexp.MustCompile(tmplNum)
)

func (c *client) Seek(key string, page int, by mart.SearchOrder) ([]mart.Product, int, error) {
	form := url.Values{
		"skeyword":  []string{key},
		"sortValue": []string{"3"},
		"offset":    []string{strconv.Itoa((page - 1) * searchSize)},
		"max":       []string{strconv.Itoa(searchSize)},
	}
	if by == mart.ByPopular {
		form["sortValue"][0] = "2"
	}

	r, err := http.Get(searchURL + "?" + form.Encode())
	if err != nil {
		return nil, 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("wellcome: Search key %s with status %s returned.", key, r.Status)
	}

	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, 0, err
	}

	// extract page number
	pg := 1
	n := regNum.FindSubmatch(o)
	if len(n) == 2 {
		pg, _ = strconv.Atoi(string(n[1]))
		pg = (pg + searchSize - 1) / searchSize
	}

	// devide html into small parts and extract the product info
	var ps []mart.Product
	frags := bytes.Split(o, []byte(`<figure class="item-image-container">`))[1:]
	for i := range frags {
		frags[i] = bytes.SplitN(frags[i], []byte(`<div class="ratings-container pull-right ">`), 2)[0]

		inf := regInfo.FindSubmatch(frags[i])
		pce := regPrice.FindSubmatch(frags[i])
		if len(inf) == 4 {
			p := mart.Product{
				Name:  string(inf[3]),
				Image: baseURL + string(inf[2]),
				Page:  baseURL + string(inf[1]),
				Mart:  c.ID(),
			}
			p.Price, _ = strconv.Atoi(string(pce[1]))
			ps = append(ps, p)
		}
	}

	// it seems Wellcome doesn't sort data completely
	// we sort it again
	if by == mart.ByPrice {
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Price-ps[j].Price < 0
		})
	}

	return ps, pg, nil
}
