package mart

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
