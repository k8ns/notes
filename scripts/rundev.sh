#!/bin/sh

docker run --rm -ti --name notes \
--network dev \
-v "$PWD":/go/src/notes -w /go/src/notes golang:1.13.3 go run ./cmd/main.go
