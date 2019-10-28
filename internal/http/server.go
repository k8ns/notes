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

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	
	r.Use(corsMiddleware)
	r.Use(gin.LoggerWithFormatter(LogFormatter))
	InitRoutes(r)
	InitWelcome(r, project)
	return r.Run(fmt.Sprintf(":%d", config.Port))
}
