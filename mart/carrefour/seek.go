package carrefour

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Hunsin/bee/mart"
)

const (
	searchURL = baseURL + "/CarrefourECProduct/GetSearchJson"

	// default number of items per search
	searchSize = 1 << 8
)

type searchProduct struct {
	ID    int    `json:"Id"`
	Name  string `json:"Name"`
	Image string `json:"PictureUrl"`
	Page  string `json:"SeName"`
	Price string `json:"Price"`

	// some products may have discount
	Special string `json:"SpecialPrice"`
}

type searchContent struct {
	Count int             `json:"Count"`
	List  []searchProduct `json:"ProductListModel"`
}

type searchJSON struct {
	Success int           `json:"success"`
	Content searchContent `json:"content"`
}

func (c *client) Seek(key string, page int, by mart.SearchOrder) ([]mart.Product, int, error) {
	form := url.Values{
		"key":       []string{key},
		"orderBy":   []string{"10"}, // by price
		"pageIndex": []string{strconv.Itoa(page)},
		"pageSize":  []string{strconv.Itoa(searchSize)},
	}
	if by == mart.ByPopular {
		form["orderBy"][0] = "21"
	}

	r, err := http.PostForm(searchURL, form)
	if err != nil {
		return nil, 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("carrefour: search key %s with status %s returned", key, r.Status)
	}

	out := searchJSON{}
	err = json.NewDecoder(r.Body).Decode(&out)
	if err != nil {
		return nil, 0, err
	}

	var ps []mart.Product
	for _, s := range out.Content.List {

		// The API may return items not match, filter it
		if !strings.ContainsAny(s.Name, key) {
			continue
		}

		p := mart.Product{
			Name:  s.Name,
			Image: s.Image,
			Page:  baseURL + strings.Split(s.Page, "?")[0],
			Mart:  c.ID(),
		}

		if s.Special != "" && s.Special != "0" {
			p.Price, err = strconv.Atoi(s.Special)
		} else {
			p.Price, err = strconv.Atoi(s.Price)
		}
		if err != nil {
			return nil, 0, err
		}

		ps = append(ps, p)
	}

	return ps, (out.Content.Count + searchSize - 1) / searchSize, nil
}
