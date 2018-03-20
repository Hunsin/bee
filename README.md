# bee [![Build Status](https://travis-ci.org/Hunsin/bee.svg?branch=master)](https://travis-ci.org/Hunsin/bee)
A crawler of online store(RT-Mart, Wellcome, Carrefour) which provides
keyword searching service through RESTful and gRPC APIs.

## Install
```sh
$ go get github.com/Hunsin/bee
```

## Usage
```sh
# Example:  
# Serve RESTful at 8888 port (default 8203)
# Serve gRPC    at 8889 port (default 8202)
$ bee -p 8888 -g 8889

# Example:
# Read version and exit
$ bee -v
```

## APIs
The API documents are under `./doc` directory.

## gRPC Client Example
An gRPC client example is under `./example` directory.

```sh
# Run client app, dials to URL "localhost:8888" (default "localhost:8202")
$ go run grpc_client.go -u localhost:8889
```
