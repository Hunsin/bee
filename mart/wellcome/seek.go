package wellcome

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"

	hp "github.com/Hunsin/bee/htmlparser"
	"github.com/Hunsin/bee/mart"
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

// count sets the number of items found to c.
func count(c *int) hp.MatchFunc {
	return func(n *html.Node) (found bool) {
		if found = hp.HasText(n, "關鍵字"); found {
			s := regNum.FindStringSubmatch(n.Data)
			if len(s) == 2 {
				*c, _ = strconv.Atoi(s[1])
			}
		}
		return
	}
}

// container locates the container of products list and assigns
// the node to c.
func container(c *html.Node) hp.MatchFunc {
	return func(n *html.Node) (found bool) {
		if found = hp.IsElement(n, "div") &&
			hp.HasAttr(n, "class", "category-item-container"); found {
			*c = *n
		}
		return
	}
}

// image fills p.Image, p.Name and p.Page by parsing attributes of
// the product's image node.
func image(p *mart.Product) hp.MatchFunc {
	return func(n *html.Node) (found bool) {
		if found = hp.IsElement(n, "img") && hp.HasAttr(n, "class", "item-image"); found {
			p.Image = baseURL + hp.Attr(n, "src")
			p.Name = hp.Attr(n, "alt")
			p.Page = baseURL + hp.Attr(n.Parent, "href")
		}
		return
	}
}

// price fills p.Price by parsing the price tag.
func price(p *mart.Product) hp.MatchFunc {
	return func(n *html.Node) (found bool) {
		if found = hp.IsElement(n, "span") && hp.HasAttr(n, "class", "item-price "); found {
			p.Price, _ = hp.Int(n)
		}
		return
	}
}

// item appends a mart.Product to ps once it found the product item node.
func item(ps *[]mart.Product) hp.MatchFunc {
	return func(n *html.Node) (found bool) {
		if found = hp.IsElement(n, "div") && hp.HasAttr(n, "class", "item"); found {
			p := mart.Product{Mart: id}
			hp.Walk(n, price(&p))
			hp.Walk(n, image(&p))

			*ps = append(*ps, p)
		}
		return
	}
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

	// extract items number, convert to page number
	var num int
	hp.Walk(n, count(&num))
	if num != 0 {
		num = (num + searchSize - 1) / searchSize
	} else {
		num = 1
	}

	// find products list container
	hp.Walk(n, container(n))

	// fill the list
	var ps []mart.Product
	hp.Walk(n, item(&ps))

	// it seems Wellcome doesn't sort data completely
	// we sort it again
	if by == mart.ByPrice {
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Price-ps[j].Price < 0
		})
	}

	return ps, num, nil
}
