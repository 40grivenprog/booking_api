package common

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	RequestIDKey string = "request_id"
	LoggerKey    string = "logger"
)

func GetRequestID(c *gin.Context) string {
	return c.GetString(string(RequestIDKey))
}

func GetLogger(c *gin.Context) zerolog.Logger {
	if logger, exists := c.Get("logger"); exists {
		return logger.(zerolog.Logger)
	}
	return zerolog.Nop()
}
