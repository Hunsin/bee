package mart

import (
	"context"
	"sync"
)

// A SearchOrder defines how data is sorted.
type SearchOrder int

const (
	// ByPrice indicates that data should sorts by price, ascending.
	ByPrice SearchOrder = iota

	// ByPopular indicates that data should sorts by popularity, descending.
	ByPopular
)

// A Query is a search request with specific keyword, how to sort
// the result and what to do after the job is done.
type Query struct {
	Key string

	// Order is either ByPopular or ByPrice.
	Order SearchOrder

	// Done is called once the search job is finished.
	// It won't be executed if the job is cancelled.
	Done func()
}

type query struct {
	ctx  context.Context
	opt  Query
	put  chan []Product
	err  chan error
	mart *Mart
	wg   sync.WaitGroup
}

// next checks if request had been cancelled, else calls fn.
func (q *query) next(fn func()) {
	select {
	case <-q.ctx.Done():
		return
	default:
		fn()
	}
}

// seek is the shorthand of q.mart.c.Seek(q.opt.Key, page, q.opt.Order)
func (q *query) seek(page int) ([]Product, int, error) {
	return q.mart.c.Seek(q.opt.Key, page, q.opt.Order)
}

// search parses the Products in given page index and sends to q.put.
// If error occurred, it will send error to q.err.
// q.Add(1) must be called before calling search.
func (q *query) search(page int) {
	defer q.wg.Done()

	// create a goroutine when first called
	// once all seek goroutines are finished, run callback
	if page == 1 && q.opt.Done != nil {
		go func() {
			q.wg.Wait()
			q.next(q.opt.Done)
		}()
	}

	// we check the channel at the beginning to avoid making request
	// after it's cancelled
	q.next(func() {
		p, m, err := q.seek(page)
		if err != nil {
			q.next(func() { q.err <- err })
			return
		}

		// if this is the first search, search the rest concurrently
		if page == 1 {
			for i := 2; i <= m; i++ {
				q.wg.Add(1)
				go q.search(i)
			}
		}

		q.next(func() { q.put <- p })
	})
}

// Search sends the slices of Product which match the given query to cp.
// If an error occurred, it sends the error to ce. It is the caller's
// responsibility to decide whether to cancel if an error is received.
func (m *Mart) Search(ctx context.Context, q Query, cp chan []Product, ce chan error) {
	qry := &query{
		ctx:  ctx,
		opt:  q,
		put:  cp,
		err:  ce,
		mart: m,
	}

	qry.wg.Add(1)
	go qry.search(1)
}
