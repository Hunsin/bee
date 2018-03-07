package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"proto"
	"rest"

	_ "carrefour"
	_ "rt"
	_ "wellcome"
)

// Version is the date of Git commit in the project.
// Rewrite it by option -ldflags="-X 'main.Version=$GIT_DATE'" in go build.
var Version = "07Mar2018"

func main() {

	// parse flags
	port := flag.Int("p", 8203, "Port of RESTful server")
	grpc := flag.Int("g", 8202, "Port of gRPC server")
	ver := flag.Bool("v", false, "Print application version and exit")
	flag.Parse()

	// print version and exit
	if *ver {
		fmt.Println("version:", Version)
		os.Exit(0)
	}

	// start gRPC server
	go func() {
		log.Println("gRPC server listen at", *grpc)
		log.Fatalln(proto.Serve(*grpc))
	}()

	// start RESTful server
	log.Println("RESTful server listen at", *port)
	log.Fatalln(rest.Serve(*port))
}
