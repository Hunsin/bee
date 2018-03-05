package mart

import (
	"errors"
	"sync"
)

var (
	pool = make(map[string]Client)
	pmu  sync.Mutex
)

// A Client is an adapter of a specific online shop.
type Client interface {

	// Seek returns the slice of Products which name match given key
	// in certain number of page. The returned integer is the number
	// of pages in total.
	Seek(string, int) ([]Product, int, error)
}

// A Product represents an item which is sold on a Mart.
type Product struct {
	Name  string // Product Name
	Image string // URL to the Product Image
	Page  string // URL of the Product page
	Price string // you know, just price
	Mart  string // The mart the product belongs to
}

// A Mart is a crawler of a online shop like RT-Mart or Carrefour.
type Mart struct {
	c Client // mart client
	n string // mart name
}

// Name returns the name of the store.
func (m *Mart) Name() string {
	return m.n
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

	return &Mart{c, name}, nil
}
