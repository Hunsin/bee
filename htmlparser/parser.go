// Package htmlparser provides some html utility functions.
package htmlparser

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// A MatchFunc returns if a html node matches certain condition.
type MatchFunc func(*html.Node) bool

// Walk calls fn with n and all the children under n. If the fn
// returns true, it stops searching the node's children.
func Walk(n *html.Node, fn MatchFunc) {
	if fn(n) {
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Walk(c, fn)
	}
}

// Attr returns the value of keyed attribute under n. If no attribute
// matches the key, an empty string is returned.
func Attr(n *html.Node, key string) string {
	for i := range n.Attr {
		if n.Attr[i].Key == key {
			return n.Attr[i].Val
		}
	}
	return ""
}

// HasAttr returns whether n has any attribute with given key
// and val.
func HasAttr(n *html.Node, key, val string) bool {
	for i := range n.Attr {
		if n.Attr[i].Key == key && n.Attr[i].Val == val {
			return true
		}
	}
	return false
}

// HasText returns wether n is a html.TextNode and n.Data contains
// given sub string.
func HasText(n *html.Node, sub string) bool {
	return n.Type == html.TextNode && strings.Contains(n.Data, sub)
}

// IsElement returns whether n is a html.ElementNode and n.Data
// matches given data.
func IsElement(n *html.Node, data string) bool {
	return n.Type == html.ElementNode && n.Data == data
}

// Text returns the content of n's first text child node. If no
// html.TextNode was found, it returns empty string.
func Text(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}
	return ""
}

// Int returns the number of n's first text child node. If no
// html.TextNode was found, it returns error.
func Int(n *html.Node) (int, error) {
	str := Text(n)
	if str == "" {
		return 0, errors.New("htmlparser: no number was found")
	}
	return strconv.Atoi(str)
}
