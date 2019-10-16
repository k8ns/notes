package main

import (
	"github.com/ksopin/notes/pkg/app"
	"github.com/ksopin/notes/pkg/db"
	"github.com/ksopin/notes/pkg/http"
)

func main() {

	cfg, err := app.ParseConfig("config.yml")
	if err != nil {
		panic(err)
	}

	db.InitConnection(cfg)

	err = http.Run(cfg)
	if err != nil {
		panic(err)
	}
}
