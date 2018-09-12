#!/bin/sh

export DATA_SOURCE="root:dbpsswrd@tcp(172.17.0.2:3306)/notes?parseTime=true"
docker run --rm -ti --env DATA_SOURCE --name notes -v "$PWD":/go/src/notes -w /go/src/notes golang:1.10.3 go run main.go