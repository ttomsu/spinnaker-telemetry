#!/usr/bin/env bash
set -x
set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
pushd "${SCRIPT_DIR}"

GOPATH=${GOPATH:=$HOME/go}
GOBIN=${GOBIN:=$GOPATH/bin}

go get -u github.com/golang/protobuf/protoc-gen-go

if [[ ":$PATH:" != *":${GOBIN}:"* ]]; then
	echo "Adding $GOBIN to PATH"
	export PATH=$PATH:$GOBIN
fi


protoc --go_out=source_relative:. *.proto
popd

