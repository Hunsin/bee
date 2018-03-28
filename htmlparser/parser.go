// Package htmlparser provides some html utility functions.
package htmlparser

import (
	"strings"

	"golang.org/x/net/html"
)

// Walk calls fn with n and all the children under n. If the fn
// returns true, it stops searching the node's children.
func Walk(n *html.Node, fn func(*html.Node) bool) {
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
