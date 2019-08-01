#!/usr/bin/env bash
set -x

GOPATH=`pwd`/.gopath
export GOPATH=$GOPATH

go get -u github.com/golang/protobuf/protoc-gen-go
export PATH=$PATH:$GOPATH/bin

./proto/regen-proto.sh

echo "gopath: $GOPATH"
#go build
