package mart

import (
	"fmt"
	"regexp"
	"strings"
)

// A Client is an adapter of a specific online store.
type Client interface {

	// Info returns the Mart's information.
	Info() Info

	// Seek returns the slice of Products which name match given key
	// in certain number of page. The third argument determines how
	// products are sorted, either ByPopular or ByPrice. The returned
	// integer is the number of pages in total.
	Seek(string, int, SearchOrder) ([]Product, int, error)
}

// ValidSeek provides an unified way to check if a client works
// as expected. It is only used for testing.
func ValidSeek(rpage, rimg string, c Client) error {
	regPage, err := regexp.Compile(rpage)
	if err != nil {
		return err
	}
	regImage, err := regexp.Compile(rimg)
	if err != nil {
		return err
	}

	for _, key := range []string{
		"抽取衛生紙",
		"蘋果",
		"牛奶花生",
		"牛排",
	} {
		ps, _, err := c.Seek(key, 1, ByPrice)
		if err != nil {
			return err
		}

		if len(ps) == 0 {
			return fmt.Errorf("search %s: no items were returned", key)
		}

		for i := range ps {
			if !regPage.MatchString(ps[i].Page) {
				return fmt.Errorf("page URL not match: %s", ps[i].Page)
			}
			if !regImage.MatchString(ps[i].Image) {
				return fmt.Errorf("image URL not match: %s", ps[i].Image)
			}
			if !strings.ContainsAny(ps[i].Name, key) {
				return fmt.Errorf("search key %s not match: %s", key, ps[i].Name)
			}

			if i == 0 {
				continue
			}

			// check if order by price
			if ps[i].Price < ps[i-1].Price {
				return fmt.Errorf("search key %s not order by price", key)
			}
		}
	}

	return nil
}
