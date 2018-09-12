package main

import (
	"github.com/gin-gonic/gin"
	"notes/http"
)

func main() {
	r := gin.Default()

	r.Use(addHeaders)
	http.InitRoutes(r)
	r.Run(":80")
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}