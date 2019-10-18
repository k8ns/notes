package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/internal/app"
)

type Config struct {
	Enabled bool
	Port int
}

func Run(config *Config, project *app.Config) error {
	if !config.Enabled {
		return nil
	}

	r := gin.Default()

	r.Use(corsMiddleware)
	InitRoutes(r)
	InitWelcome(r, project)
	return r.Run(fmt.Sprintf(":%d", config.Port))
}
