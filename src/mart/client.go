package mart

import (
	"errors"
	"regexp"
	"strings"
)

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

	key := "抽取衛生紙"
	ps, _, err := c.Seek(key, 1, ByPrice)
	if err != nil {
		return err
	}

	for i := range ps {
		if !regPage.MatchString(ps[i].Page) {
			return errors.New("page not match: " + ps[i].Page)
		}
		if !regImage.MatchString(ps[i].Image) {
			return errors.New("image not match: " + ps[i].Image)
		}
		if !strings.ContainsAny(ps[i].Name, key) {
			return errors.New("name not match: " + ps[i].Name)
		}

		if i == 0 {
			continue
		}

		// check if order by price
		if ps[i].Price < ps[i-1].Price {
			return errors.New("not order by price")
		}
	}

	return nil
}
