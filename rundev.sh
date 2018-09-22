#!/bin/sh

export DATA_SOURCE="root:dbpsswrd@tcp(172.18.0.3:3306)/notes?parseTime=true"
docker run --rm -ti --env DATA_SOURCE --name notes \
--network dev --ip 172.18.0.30 \
-v "$PWD":/go/src/notes -w /go/src/notes golang:1.11 go run main.go
