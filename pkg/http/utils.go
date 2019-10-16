package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func writeOkResponse(c *gin.Context, data interface{}, status int) {
	c.JSON(status, gin.H{
		"data": data,
	})
}

func writeErrResponse(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func writeMapErrResponse(c *gin.Context, errs map[string]error, status int) {
	m := make(map[string]string, len(errs))
	for key, err := range errs {
		m[key] = err.Error()
	}
	c.JSON(status, gin.H{
		"error": m,
	})
}


func ok(c *gin.Context) {
	c.Status(http.StatusOK)
}
