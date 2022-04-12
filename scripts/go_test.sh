#!/bin/bash

source ~/.bashrc
gvm install go1.16
gvm use go1.16

# go mod download
go mod tidy

/usr/local/bin/gocheck -t
