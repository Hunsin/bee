package mart

import (
	"fmt"
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
