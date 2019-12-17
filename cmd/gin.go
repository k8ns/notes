package main

import (
	"github.com/ksopin/notes/internal/http"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	err := http.Run()
	if err != nil {
		panic(err)
	}
}
