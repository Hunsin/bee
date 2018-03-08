// Package mart is the core part of the model. It extracts products
// information from online store through a registered adapter.
// The package must be used in conjunction with adapters.
package mart

import (
	"errors"
	"sync"
)

// Currencies
const (
	CurrencyTWD = "TWD"
	CurrencyUSE = "USD"
)

var (
	pool = make(map[string]Client)
	pmu  sync.Mutex
)

// A Client is an adapter of a specific online store.
type Client interface {

	// Currency returns the currency the Mart is use.
	Currency() string

	// ID returns the abbreviation of the Mart.
	ID() string

	// Name returns the full name of the Mart.
	Name() string

	// Seek returns the slice of Products which name match given key
	// in certain number of page. The third argument determines how
	// products are sorted, either ByPopular or ByPrice. The returned
	// integer is the number of pages in total.
	Seek(string, int, SearchOrder) ([]Product, int, error)
}

// A Product represents an item which is sold on a Mart.
type Product struct {
	Name  string `json:"name"`  // Product Name
	Image string `json:"image"` // URL to the Product Image
	Page  string `json:"page"`  // URL of the Product page
	Price int    `json:"price"` // you know, just price
	Mart  string `json:"mart"`  // The mart the product belongs to
}

// A Info specifies the information of a Mart
type Info struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"cur"`
}

// A Mart is a crawler of a online store like RT-Mart or Carrefour.
type Mart struct {
	c Client // mart client
}

// Info returns the information of the store.
func (m *Mart) Info() Info {
	return Info{
		m.c.ID(),
		m.c.Name(),
		m.c.Currency(),
	}
}

// Register makes a client available by the provided name. If c is nil
// or Register is called twice with the same name, it panics.
func Register(name string, c Client) {
	pmu.Lock()
	defer pmu.Unlock()

	if c == nil {
		panic("mart: A nil Client is registered")
	}

	if _, ok := pool[name]; ok {
		panic("mart: Multiple Clients registered under name " + name)
	}

	pool[name] = c
}

// Open returns a pointer to Mart with named Client.
func Open(name string) (*Mart, error) {
	c, ok := pool[name]
	if !ok {
		return nil, errors.New("Client " + name + " not found")
	}

	return &Mart{c}, nil
}

// All returns all Marts available.
func All() []*Mart {
	var m []*Mart
	for _, c := range pool {
		m = append(m, &Mart{c})
	}

	return m
}
