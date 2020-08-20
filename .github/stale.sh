#!/bin/bash -ex

export PATH=$PATH:$(pwd)/bin
export GO111MODULE=on
export GOBIN=$(pwd)/bin

#go get github.com/rvflash/goup@v0.4.1

#goup -v ./...
#go get github.com/psampaz/go-mod-outdated@v0.6.0
go list -u -m -mod=mod -json all | go-mod-outdated -update -direct -ci || true

#go list -u -m -json all | go-mod-outdated -update
