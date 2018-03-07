# bee
A crawler of online store of marts which provides RESTful and gRPC APIs.

## Build
1. Clone it
```sh
$ git clone https://github.com/Hunsin/bee.git
$ cd ./bee
```
2. Build it
```sh
$ ./build.sh all
```
3. Done! Do you need more information?
```sh
$ ./build.sh help
```

## Run
Once you build it, the app is under `./dist` directory.

```sh
# Example:  
# Serve RESTful at 8888 port (default 8203)
# Serve gRPC    at 8889 port (default 8202)
$ cd ./dist/
$ ./app -p 8888 -g 8889

# Example:
# Read version and exit
$ ./app -v
```

## APIs
The API documents are under `./doc` directory.

## gRPC Client
```sh
# Run client app, dials to URL "localhost:8888"
$ export GOPATH=`pwd`
$ go run dist/client.go -u localhost:8889
```
