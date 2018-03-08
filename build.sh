#!/bin/bash
set -e

BUILD_PATH=$(cd $(dirname $0); pwd)
DIST_PATH=$BUILD_PATH/dist
GIT_DATE=$(git log -1 --format=%cd)

PROTOC_ZIP=protoc-3.5.1-linux-x86_64.zip
PROTOC_URL=https://github.com/google/protobuf/releases/download/v3.5.1/$PROTOC_ZIP

export GOPATH=$BUILD_PATH
export PATH=$PATH:$GOPATH/bin

# env installs necessary packages
function env(){

	# install packages
	echo "Start installing golang packages..."
	echo "It take some time, please wait..."
  go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc

	# download gRPC compiler
	echo "Start downloading protobuf compiler..."
	mkdir $BUILD_PATH/protoc
	wget $PROTOC_URL -O $BUILD_PATH/protoc/$PROTOC_ZIP
	unzip $BUILD_PATH/protoc/$PROTOC_ZIP -d $BUILD_PATH/protoc

	# make directory where files are output
	mkdir $BUILD_PATH/src/proto/pb
}

# app builds the crawler which provides RESTful and gRPC APIs for
# clients to scrape data from online store
function app(){
	cd $BUILD_PATH
	protoc/bin/protoc -I doc/ doc/APIs_client.proto --go_out=plugins=grpc:src/proto/pb

	go build -o $DIST_PATH/app -ldflags="-X 'main.Version=$GIT_DATE'" main
	echo "Binary file saved as: $DIST_PATH/app"
}

# all calls env() and app()
function all(){
	env
	app
}

function help(){
cat <<EOF
Usage:
$ ./build.sh app: Build application.
$ ./build.sh env: Install necessary packages.
$ ./build.sh all: Execute "./build.sh env" and "./build.sh app".
EOF
}

case $1 in
	app) app;;
	env) env;;
	all) all;;
	*)   help;;
esac
