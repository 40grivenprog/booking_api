package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeAuth, common.ErrorMsgMissingAuthToken, nil)
			c.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeAuth, common.ErrorMsgInvalidAuthHeader, nil)
			c.Abort()
			return
		}

		authorizationType := fields[0]
		if authorizationType != authorizationTypeBearer {
			common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeAuth, common.ErrorMsgUnsupportedAuthType, nil)
			c.Abort()
			return
		}

		accessToken := fields[1]
		_, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeAuth, common.ErrorMsgInvalidToken, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
