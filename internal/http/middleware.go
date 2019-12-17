package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/pkg/auth"
)

func CorsMiddleware (c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func authMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatus(401)
		return
	}

	u, err := app.GetAuthService().VerifySignature(c, token)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	if u == nil {
		c.AbortWithStatus(403)
		return
	}

	c.Set(auth.AuthUserKey, u)
}
