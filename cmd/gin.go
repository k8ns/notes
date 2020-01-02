package main

import (
	"github.com/ksopin/notes/internal/config"
	"github.com/ksopin/notes/internal/http"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	cfg := config.GetConfig(strings.Join([]string{"config/config", os.Getenv("APP_ENV"), "yml"}, "."))

	log.SetFormatter(&log.JSONFormatter{})
	err := http.Run(cfg.Http, cfg.App)
	if err != nil {
		panic(err)
	}
}
