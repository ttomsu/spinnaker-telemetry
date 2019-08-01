#!/usr/bin/env bash
set -x
set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
pushd "${SCRIPT_DIR}"

mkdir -p ${GOPATH:=$HOME/go}/src
./protoc log.proto --go_out=${GOPATH:=$HOME/go}/src
popd

