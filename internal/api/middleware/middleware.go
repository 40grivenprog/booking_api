package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/vention/booking_api/internal/api/common"
)

const RequestIDKey = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Writer.Header().Set(RequestIDKey, requestID)

		c.Set(common.RequestIDKey, requestID)

		c.Next()
	}
}

// Logger creates an adjusted logger with request context and logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get the base logger from context (set by server)
		baseLogger := c.MustGet("logger").(zerolog.Logger)

		// Create adjusted logger with request context
		adjustedLogger := baseLogger.With().
			Str("request_id", common.GetRequestID(c)).
			Str("endpoint", c.FullPath()).
			Str("method", c.Request.Method).
			Logger()

		// Set adjusted logger in context for use in handlers
		c.Set(common.LoggerKey, adjustedLogger)

		// Log request start
		adjustedLogger.Info().Msg("Request started")

		c.Next()

		// Log request completion with timing and status
		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		adjustedLogger.Info().
			Int("status", statusCode).
			Int("size", bodySize).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Msg("Request completed")
	}
}
