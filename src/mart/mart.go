package mart

import (
	"errors"
	"sync"
)

var (
	pool = make(map[string]Driver)
	pmu  sync.Mutex
)

type Driver interface {
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
	d Driver // mart driver
	n string // mart name
}

// Name returns the name of the store.
func (m *mart) Name() string {
	return m.n
}

func Register(name string, d Driver) {
	pmu.Lock()
	defer pmu.Unlock()

	if d == nil {
		panic("mart: A nil Driver is registered")
	}

	if _, ok := pool[name]; ok {
		panic("mart: Multiple Drivers registered under name " + name)
	}

	pool[name] = d
}

func Open(name string) (*Mart, error) {
	d, ok := pool[name]
	if !ok {
		return nil, errors.New("Driver " + name + " not found")
	}

	return d, nil
}
