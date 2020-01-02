package main

import (
	"github.com/ksopin/notes/internal/config"
	ginlambda "github.com/ksopin/notes/internal/lambda"
	"os"
	"strings"
)

func main() {
	cfg := config.GetConfig(strings.Join([]string{"config/config", os.Getenv("APP_ENV"), "yml"}, "."))

	err := ginlambda.Run(cfg.App)
	if err != nil {
		panic(err)
	}
}
