package mart

import "sync"

type searchReq struct {
	key  string
	put  chan []Product
	err  chan error
	done chan bool
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

func (s *searchReq) seek(page int) {
	s.wg.Add(1)
	defer s.wg.Done()

	// we check the channel at the beginning to avoid making request
	// after it's cancelled
	s.next(func() {
		p, m, err := s.mart.d.Seek(s.key, page)
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
				s.next(func() { s.done <- true })
			}()
		}

		s.next(func() { s.put <- ps })
	})
}

// Search returns a channel of Products which name match given key
// and a channel of error. It is the caller's responsibility to
// decide whether to keep going if an error is received. Closing the
// done notifies the Mart to stop searching. If no more Product left,
// it will send true to done.
func (m *Mart) Search(key string, done chan bool) (chan []Product, chan error) {
	s := &searchReq{
		key:  key,
		put:  make(chan []Product),
		err:  make(chan error),
		done: done,
		mart: m,
	}

	go s.search(1, done)
	return s.put, s.err
}
