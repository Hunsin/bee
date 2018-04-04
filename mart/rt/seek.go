package rt

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/Hunsin/bee/mart"
	hu "github.com/Hunsin/go-htmlutil"
	"golang.org/x/net/html"
)

var numSeek = 1 << 6

// count returns the number of items found.
// The number is extracted from "span.t02" element.
func count(doc *html.Node) (c int) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.IsElement(n, "span") && hu.HasAttr(n, "class", "t02"); found {
			var err error
			if c, err = hu.Int(n); err != nil {
				c = 1
			}
		}
		return
	})
	return
}

// container returns the pointer to the node of products list.
func container(doc *html.Node) (c *html.Node) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.HasAttr(n, "class", "classify_prolistBox"); found {
			c = n
		}
		return
	})
	return
}

// image extracts the product name, image and page URL from
// "div.for_imgbox" element. The p.Price is 0.
func image(doc *html.Node) (p *mart.Product) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.IsElement(n, "img"); found {
			p = &mart.Product{
				Image: hu.Attr(n, "src"),
				Name:  hu.Attr(n, "title"),
				Page:  hu.Attr(n.Parent, "href"),
				Mart:  id,
			}
		}
		return
	})
	return
}

// price returns the product price.
func price(doc *html.Node) (p int) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.HasAttr(n, "class", "for_pricebox"); found {
			p, _ = hu.Int(n)
		}
		return
	})
	return
}

// list returns the slice of products extracted from the page.
func list(doc *html.Node) (ps []mart.Product) {
	hu.Walk(doc, func(n *html.Node) (found bool) {
		if found = hu.HasAttr(n, "class", "indexProList"); found {
			if p := image(n); p != nil {
				p.Price = price(n)
				ps = append(ps, *p)
			}
		}
		return
	})
	return
}

func (c *client) Seek(key string, page int, by mart.SearchOrder) ([]mart.Product, int, error) {
	form := url.Values{
		"action":       []string{"product_search"},
		"prod_keyword": []string{key},
		"p_data_num":   []string{strconv.Itoa(numSeek)},
		"page":         []string{strconv.Itoa(page)},
		"usort":        []string{"prod_selling_price,ASC"},
	}
	if by == mart.ByPopular {
		form["usort"][0] = "prod_sales_count,DESC"
	}

	r, err := http.Get(baseURL + "?" + form.Encode())
	if err != nil {
		return nil, 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("rt: search key %s with status %s returned", key, r.Status)
	}

	n, err := html.Parse(r.Body)
	if err != nil {
		return nil, 0, err
	}

	// assign pointer to <body></body> element
	hu.First(n, func(c *html.Node) (found bool) {
		if found = hu.IsElement(c, "body"); found {
			*n = *c
		}
		return
	})

	// extract number of products
	num := count(n)

	// get product list node
	n = container(n)

	ps := list(n)

	// it seems RT-Mart doesn't sort data completely
	// we sort it again
	if by == mart.ByPrice {
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Price-ps[j].Price < 0
		})
	}

	return ps, (num + numSeek - 1) / numSeek, nil
}
