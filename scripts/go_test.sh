#!/bin/bash

source ~/.bashrc
gvm install $1
gvm use $1
go mod tidy

# go test
go test /go/src/github.com/y-miyazaki/go-common/...
