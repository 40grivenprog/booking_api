package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	RequestIDKey string = "request_id"
	LoggerKey    string = "logger"
	UserIDKey    string = "user_id"
	UserNameKey  string = "user_name"
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

// GetUserID retrieves the user ID from the context (set by AuthMiddleware)
// Returns the user ID and a boolean indicating success
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDValue, exists := c.Get(UserIDKey)
	if !exists {
		HandleErrorResponse(c, http.StatusUnauthorized, ErrorTypeAuth, ErrorMsgMissingAuthToken, nil)
		return uuid.UUID{}, false
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		HandleErrorResponse(c, http.StatusInternalServerError, ErrorTypeInternal, "Invalid user ID type", nil)
		return uuid.UUID{}, false
	}

	return userID, true
}

func GetUserName(c *gin.Context) string {
	userNameValue, exists := c.Get(UserNameKey)
	if !exists {
		HandleErrorResponse(c, http.StatusUnauthorized, ErrorTypeAuth, ErrorMsgMissingAuthToken, nil)
		return ""
	}
	return userNameValue.(string)
}
