package wellcome

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"

	"github.com/Hunsin/bee/mart"
	hu "github.com/Hunsin/go-htmlutil"
	"golang.org/x/net/html"
)

const (
	searchURL = baseURL + "/product/listByKeyword"

	// max numer of items per search
	searchSize = 100

	// number of products
	tmplNum = `.+\(([0-9]*)\)`
)

var regNum = regexp.MustCompile(tmplNum)

// count returns the number of items found.
func count(doc *html.Node) (c int) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.HasText(n, "關鍵字"); found {
			s := regNum.FindStringSubmatch(n.Data)
			if len(s) == 2 {
				c, _ = strconv.Atoi(s[1])
			}
		}
		return
	})
	return
}

// container returns the pointer to the node of products container,
// which the element is "div.category-item-container"
func container(doc *html.Node) (c *html.Node) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.IsElement(n, "div") &&
			hu.HasAttr(n, "class", "category-item-container"); found {
			c = n
		}
		return
	})
	return
}

// image returns the product name, image and page URL from
// "img.item-image" element. The p.Mart is set to id.
func image(doc *html.Node) (p *mart.Product) {
	hu.First(doc, func(n *html.Node) (found bool) {
		if found = hu.IsElement(n, "img") && hu.HasAttr(n, "class", "item-image"); found {
			p = &mart.Product{
				Image: baseURL + hu.Attr(n, "src"),
				Name:  hu.Attr(n, "alt"),
				Page:  baseURL + hu.Attr(n.Parent, "href"),
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
		if found = hu.IsElement(n, "span") && hu.HasAttr(n, "class", "item-price "); found {
			p, _ = hu.Int(n)
		}
		return
	})
	return
}

// list returns the slice of products extracted from the page.
func list(doc *html.Node) (ps []mart.Product) {
	hu.Walk(doc, func(n *html.Node) (found bool) {
		if found = hu.IsElement(n, "div") && hu.HasAttr(n, "class", "item"); found {
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
		return nil, 0, fmt.Errorf("wellcome: search key %s with status %s returned", key, r.Status)
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

	// extract number of items
	num := count(n)

	// find products list container
	n = container(n)

	// fill the list
	ps := list(n)

	// it seems Wellcome doesn't sort data completely
	// we sort it again
	if by == mart.ByPrice {
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Price-ps[j].Price < 0
		})
	}

	return ps, (num + searchSize - 1) / searchSize, nil
}
