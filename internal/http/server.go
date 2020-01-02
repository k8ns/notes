package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/internal/app"
)

func New(appCfg *app.Config) *gin.Engine {
	app.InitApp(appCfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(CorsMiddleware)
	r.Use(gin.LoggerWithFormatter(LogFormatter))
	InitRoutes(r)
	InitWelcome(r, appCfg)

	return r
}

func Run(cfg *Config, appCfg *app.Config) error {
	if !cfg.Enabled {
		return nil
	}

	r := New(appCfg)
	return r.Run(fmt.Sprintf(":%d", cfg.Port))
}
