#!/bin/bash

source ~/.bashrc
gvm install go1.19
gvm use go1.19

go mod tidy

# go test
go test /go/src/github.com/y-miyazaki/go-common/...
