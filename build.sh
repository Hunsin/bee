#! /bin/sh
BUILD_PATH=$(cd $(dirname $0); pwd)
DIST_PATH=$BUILD_PATH/dist
GIT_DATE=$(git log -1 --format=%cd)

PROTOC_ZIP=protoc-3.5.1-linux-x86_64.zip
PROTOC_URL=https://github.com/google/protobuf/releases/download/v3.5.1/$PROTOC_ZIP

export GOPATH=$BUILD_PATH
export PATH=$PATH:$GOPATH/bin

# env installs necessary packages
function env(){
  # not sure if needs
	# go get -u github.com/Hunsin/beaver
	
  go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc

	mkdir $BUILD_PATH/protoc
	wget $PROTOC_URL -O $BUILD_PATH/protoc/$PROTOC_ZIP
	unzip $BUILD_PATH/protoc/$PROTOC_ZIP -d $BUILD_PATH/protoc
}

# app builds the crawler which provides gRPC APIs for clients
# to scrape data from online shop
function app(){
	env
	cd $BUILD_PATH
	mkdir $BUILD_PATH/src/proto/pb
	protoc/bin/protoc -I doc/ doc/APIs_client.proto --go_out=plugins=grpc:src/proto/pb

	go build -o $DIST_PATH/crawler -ldflags="-X 'main.Version=$GIT_DATE'" $BUILD_PATH/app/crawler.go
	echo "Binary file saved as: $DIST_PATH/crawler"
}

# web builds the web server which provides RESTful APIs
# and simple front-end page
function web(){
	env
  cd $BUILD_PATH
	
	go build -o $DIST_PATH/web -ldflags="-X 'main.Version=$GIT_DATE'" $BUILD_PATH/app/web.go
	echo "Binary file saved as: $DIST_PATH/web"
}

function help(){
cat <<EOF
Usage:
$ ./build.sh app: Build crawler server.
$ ./build.sh env: Install necessary packages.
$ ./build.sh web: Build web server.
EOF
}


case $1 in
	app) app;;
	env) env;;
	web) web;;
	*)   help;;
esac
