package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func LogFormatter(p gin.LogFormatterParams) string {

	b, err := json.Marshal(map[string]interface{}{
		"timeStamp": p.TimeStamp,
		"statusCode": p.StatusCode,
		"latency": p.Latency,
		"clientIP": p.ClientIP,
		"method": p.Method,
		"path": p.Path,
		"errorMessage": p.ErrorMessage,
		"bodySize": p.BodySize,
		"keys": p.Keys,
	})

	if err != nil {
		return err.Error()
	}

	return string(b) + "\n"
}
