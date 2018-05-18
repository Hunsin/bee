package server

import (
	"log"
	"net"
	"strconv"

	"github.com/Hunsin/bee/api"
	"github.com/Hunsin/bee/mart"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// badRequest returns an status error with codes.InvalidArgument.
func badRequest(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// noFound returns an status error with codes.NotFound.
func noFound(msg string) error {
	return status.Error(codes.NotFound, msg)
}

// A gRPCsrv implements the api.CrawlerServer interface.
type gRPCsrv struct{}

// Search streams the products which match q to client.
func (s *gRPCsrv) Search(q *api.Query, stream api.Crawler_SearchServer) error {
	if q.Key == "" {
		return badRequest("Key must not be empty")
	}

	// create query
	d := make(chan bool)
	opt := mart.Query{
		Key:   q.Key,
		Order: mart.ByPrice,
		Done:  func() { d <- true },
	}
	if q.Order == api.Query_POPULAR {
		opt.Order = mart.ByPopular
	}

	// find if mart available
	var ms []*mart.Mart
	if q.Mart != "" {
		m, err := mart.Open(q.Mart)
		if err != nil {
			return noFound("Mart " + q.Mart + " not available")
		}

		ms = append(ms, m)
	} else {
		ms = mart.All()
		if len(ms) == 0 {
			return noFound("No mart available")
		}
	}

	// create context and channel; make search request
	ctx, quit := context.WithCancel(stream.Context())
	defer quit()

	put := make(chan []mart.Product)
	che := make(chan error)
	for i := range ms {
		ms[i].Search(ctx, opt, put, che)
	}

	// listen for search response
	var sent, done int64
	for {
		select {
		case <-ctx.Done():
			log.Println("Search keyword", q.Key, "cancelled")
			return nil
		case ps := <-put:
			for i := range ps {
				sent++
				if q.Num > 0 && sent > q.Num { // reach max number, return
					return nil
				}

				if err := stream.Send(&api.Product{
					Name:  ps[i].Name,
					Image: ps[i].Image,
					Page:  ps[i].Page,
					Price: int64(ps[i].Price),
					Mart:  ps[i].Mart,
				}); err != nil {
					log.Println(err)
					return nil // connection lost?
				}
			}
		case err := <-che:
			log.Println(err)
		case <-d:
			done++
			if done == int64(len(ms)) { // all jobs are done
				return nil
			}
		}
	}
}

// Marts responses with the client a list of marts available.
func (s *gRPCsrv) Marts(_ *api.Null, stream api.Crawler_MartsServer) error {
	all := mart.All()
	if len(all) == 0 {
		return noFound("No mart available")
	}

	for _, m := range all {
		info := m.Info()
		stream.Send(&api.Mart{
			Id:   info.ID,
			Name: info.Name,
			Cur:  info.Currency,
		})
	}

	return nil
}

// GRPC creates a gRPC server which listens to given port.
func GRPC(port int) error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	api.RegisterCrawlerServer(s, &gRPCsrv{})

	return s.Serve(l)
}
