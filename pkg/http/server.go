package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/pkg/app"
)

func Run(config *app.Config) error {
	r := gin.Default()

	r.Use(corsMiddleware)
	InitRoutes(r)
	return r.Run(":" + config.Httpport)
}
