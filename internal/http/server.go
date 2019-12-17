package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/pkg/db"
	"os"
	"strings"
)



func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(CorsMiddleware)
	r.Use(gin.LoggerWithFormatter(LogFormatter))
	InitRoutes(r)

	return r
}

func Run() error {
	cfg := app.GetConfig(strings.Join([]string{"config/config", os.Getenv("APP_ENV"), "yml"}, "."))

	if !cfg.Http.Enabled {
		return nil
	}
	db.InitConnection(cfg.Db)
	r := New()
	InitWelcome(r, cfg.App)
	return r.Run(fmt.Sprintf(":%d", cfg.Http.Port))
}
