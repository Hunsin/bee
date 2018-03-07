package mart

import (
	"context"
	"sync"
)

// A SearchOrder defines how data is sorted.
type SearchOrder int

const (
	ByPrice SearchOrder = iota
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

// seek is the shorthand of q.mart.c.Seek(q.opt.Key, q.opt.Order, page)
func (q *query) seek(page int) ([]Product, int, error) {
	return q.mart.c.Seek(q.opt.Key, page, q.opt.Order)
}

// search parses the Products in given page index and sends to q.put.
// If error occurred, it will send error to q.err.
func (q *query) search(page int) {
	q.wg.Add(1)
	defer q.wg.Done()

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
			go func() {
				for i := 2; i <= m; i++ {
					q.search(i)
				}

				// once all seek goroutines are finished, run callback
				q.wg.Wait()
				if q.opt.Done != nil {
					q.next(q.opt.Done)
				}
			}()
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

	go qry.search(1)
}
