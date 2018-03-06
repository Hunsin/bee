package proto

import (
	"log"
	"mart"
	"net"
	"proto/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// noFound returns an status error with code NotFound.
func noFound(msg string) error {
	return status.Error(codes.NotFound, msg)
}

// A gRPCsrv implements the pb.CrawlerServer interface.
type gRPCsrv struct{}

func (s *gRPCsrv) Search(q *pb.Query, stream pb.Crawler_SearchServer) error {

	// create query
	d := make(chan bool)
	opt := mart.Query{
		Key:   q.Key,
		Order: mart.ByPrice,
		Done:  func() { d <- true },
	}
	if q.Order == pb.Query_POPULAR {
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
			return nil
		case ps := <-put:
			for i := range ps {
				sent++
				if q.Num > 0 && sent > q.Num { // reach max number, return
					return nil
				}

				if err := stream.Send(&pb.Product{
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

func (s *gRPCsrv) Marts(_ context.Context, _ *pb.Null) (*pb.MartList, error) {
	l := &pb.MartList{}
	for _, m := range mart.All() {
		info := m.Info()
		l.Marts = append(l.Marts, &pb.Mart{
			Id:   info.ID,
			Name: info.Name,
			Cur:  info.Currency,
		})
	}

	if len(l.Marts) == 0 {
		return nil, noFound("No mart available")
	}
	return l, nil
}

// Serve creates a gRPC server which listens to given port.
func Serve(port string) error {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterCrawlerServer(s, &gRPCsrv{})

	return s.Serve(l)
}
