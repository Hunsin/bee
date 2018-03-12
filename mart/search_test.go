package mart

import (
	"context"
	"errors"
	"testing"
	"time"
)

const (
	shouldError = iota - 2
	shouldCancel
)

// A mock implements the Client interface.
type mock struct {
	status int
	cancel func()
}

func (m *mock) Currency() string {
	return CurrencyUSD
}

func (m *mock) ID() string {
	return "mock"
}

func (m *mock) Name() string {
	return "Mock Client"
}

// Seek returns with page == 5. When page == 3, it calls m.cancel() 
// if m.status == shouldCancel; or it returns an error if
// m.status == shouldError.
func (m *mock) Seek(_ string, page int, odr SearchOrder) ([]Product, int, error) {
	if page == 3 {
		if m.status == shouldCancel {
			m.cancel()
		}

		if m.status == shouldError {
			return nil, 0, errors.New("error message")
		}
	}

	return []Product{Product{}, Product{}}, 5, nil
}

func TestSearch(t *testing.T) {
	done := make(chan bool)
	qry := Query{"", ByPopular, func() { done <- true }}
	cp := make(chan []Product, 10)
	ce := make(chan error)
	mk := &mock{}
	Register(mk)
	mt, _ := Open(mk.ID())

	// test error
	ctx, quit := context.WithTimeout(context.Background(), 3*time.Second)
	mk.status = shouldError
	mk.cancel = quit
	mt.Search(ctx, qry, cp, ce)

	select { // should receive error before ctx.Done()
	case err := <-ce:
		if err == nil {
			t.Error("Mart.Search failed: error not received")
		}
	case <-ctx.Done():
		t.Error("Mart.Search failed: error not received")
	}

	select { // should receive done before ctx.Done()
	case <-done:
		// continue
	case <-ctx.Done():
		t.Error("Mart.Search failed: Done func not executed")
	}

	// 4 slice should in channel buffer (another was sent to ce)
	if len(cp) != 4 {
		t.Error("Mart.Search failed: Product slice not sent. Got:", len(cp))
	}

	// test cancel
	ctx, quit = context.WithCancel(context.Background())
	mk.status = shouldCancel
	mk.cancel = quit
	mt.Search(ctx, qry, cp, ce)

	select { // should receive ctx.Done() only
	case <-done:
		t.Error("Mart.Search failed: Query.Done was called in cancelled request")
	case <-ctx.Done():
		// continue
	}
}
