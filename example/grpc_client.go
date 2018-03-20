package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	
	"github.com/Hunsin/bee/server/pb"
	"google.golang.org/grpc"
)

// send reads STDIN and sends the request to server.
func send(c pb.CrawlerClient) {
	fmt.Print("Search: ")
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {

		// avoid empty input
		in := s.Text()
		if len(in) == 0 {
			fmt.Print("Search: ")
			continue
		}

		// get keyword and max
		var max int
		var err error
		args := strings.Fields(in)
		if len(args) > 1 {
			max, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Input", args[1], "is not a valid integer")
				break
			}
		}
		if args[0] == "\\q" { // quit
			os.Exit(0)
		}

		// send request
		r, err := c.Search(context.Background(), &pb.Query{
			Key: args[0],
			Num: int64(max),
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		recv(r)
		fmt.Print("Search: ")
	}
}

// recv receives the data from server and prints the products to STDOUT.
func recv(c pb.Crawler_SearchClient) {
	var count int
	for {
		p, err := c.Recv()
		if err == io.EOF {
			fmt.Println("Done, received:", count)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		count++
		fmt.Printf("Name:  %s\nPrice: %d\nFrom:  %s\nPage:  %s\n\n", p.Name, p.Price, p.Mart, p.Page)
	}
}

func main() {
	url := flag.String("u", "localhost:8202", "Server URL")
	flag.Parse()

	c, err := grpc.Dial(*url, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := pb.NewCrawlerClient(c)
	fmt.Println("gRPC server:", *url, "connected")
	fmt.Println("Usage:   keyword    max(optional)")
	fmt.Println("Example: 抽取衛生紙 30")
	fmt.Println("To exit: \\q")
	send(client)
}
