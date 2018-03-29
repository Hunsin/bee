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
	CurrencyUSD = "USD"
)

var (
	pool = make(map[string]Client)
	pmu  sync.Mutex
)

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
	return m.c.Info()
}

// Register makes a client available for use. If c is nil or Register
// is called twice with the same client ID, it panics.
func Register(c Client) {
	pmu.Lock()
	defer pmu.Unlock()

	if c == nil {
		panic("mart: A nil Client is registered")
	}

	id := c.Info().ID
	if _, ok := pool[id]; ok {
		panic("mart: Multiple Clients registered under ID " + id)
	}

	pool[id] = c
}

// Open returns a pointer to Mart with given id.
func Open(id string) (*Mart, error) {
	c, ok := pool[id]
	if !ok {
		return nil, errors.New("Client " + id + " not found")
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
