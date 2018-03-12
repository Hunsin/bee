package rt

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"

	"github.com/Hunsin/bee/mart"
)

const (
	tmplPro = `<a href="(.+)" target="_blank"><img src="(.+)" width="[0-9]+" (?:[a-z]+=".+" )+alt="(.+)" ?/>.+<div class="for_pricebox"><div ?><span>\$</span>([0-9]+)</div>`
)

var (

	// regPro parses the product information
	regPro = regexp.MustCompile(tmplPro)

	// regNum parses the number of products
	regNum = regexp.MustCompile(`<span class="t02">([0-9]*)</span>`)

	numSeek = 1 << 6
)

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

	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, 0, err
	}

	// extract number of products
	var num int
	b := regNum.FindSubmatch(o)
	if len(b) == 2 {
		num, _ = strconv.Atoi(string(b[1]))
	}

	var ps []mart.Product

	// devide the document into small fragments and remove "\n"
	// than extract the products
	frags := bytes.Split(o, []byte(`<div class="indexProList">`))[1:]
	for i := range frags {
		frags[i] = bytes.Replace(frags[i], []byte("\n"), []byte{}, -1)

		b = regPro.FindSubmatch(frags[i])
		if len(b) == 5 {
			p := mart.Product{
				Name:  string(b[3]),
				Image: string(b[2]),
				Page:  string(b[1]),
				Mart:  c.ID(),
			}
			p.Price, _ = strconv.Atoi(string(b[4]))
			ps = append(ps, p)
		}
	}

	// it seems RT-Mart doesn't sort data completely
	// we sort it again
	if by == mart.ByPrice {
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].Price-ps[j].Price < 0
		})
	}

	return ps, (num + numSeek - 1) / numSeek, nil
}