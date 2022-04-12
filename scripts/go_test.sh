#!/bin/bash

gvm install go1.16
gvm use go1.16

# go mod download
go mod download

/usr/local/bin/gocheck -t
