package mart

import "sync"

// A searchReq represents an search event with specific keyword.
type searchReq struct {
	key  string
	put  chan []Product
	err  chan error
	done chan string
	mart *Mart
	wg   sync.WaitGroup
}

// next checks if request had been cancelled, else calls fn.
func (s *searchReq) next(fn func()) {
	select {
	case <-s.done:
		return
	default:
		fn()
	}
}

// seek parses the Products in given page index and sends to s.put.
// If error occurred, it will send error to s.err.
func (s *searchReq) seek(page int) {
	s.wg.Add(1)
	defer s.wg.Done()

	// we check the channel at the beginning to avoid making request
	// after it's cancelled
	s.next(func() {
		p, m, err := s.mart.c.Seek(s.key, page)
		if err != nil {
			s.next(func() { s.err <- err })
			return
		}

		// if this is the first search, search the rest concurrently
		if page == 1 {
			go func() {
				for i := 2; i <= m; i++ {
					s.seek(i)
				}

				// once all seek goroutines are finished, notifies the caller
				s.wg.Wait()
				s.next(func() { s.done <- s.mart.Name() })
			}()
		}

		s.next(func() { s.put <- p })
	})
}

// Search sends the slices of Product which name match the given key
// to cp. If an error occurred, it sends the error to ce. It is the
// caller's responsibility to decide whether to continue if an error
// is received. Closing done notifies the Mart to stop searching.
// Once the job is finished, it sends m.Name() to done.
func (m *Mart) Search(key string, done chan string, cp chan []Product, ce chan error) {
	s := &searchReq{
		key:  key,
		put:  cp,
		err:  ce,
		done: done,
		mart: m,
	}

	go s.seek(1)
}
