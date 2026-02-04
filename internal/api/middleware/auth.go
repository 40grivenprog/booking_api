package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
)

var skipPaths = []string{"/api/users/:chat_id", "/api/clients/register", "/api/professionals/sign_in", "/api/admins/professionals"}

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(skipPaths, c.FullPath()) {
			fmt.Println("skipping path", c.Request.URL.Path)
			c.Next()
			return
		}

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
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeAuth, common.ErrorMsgInvalidToken, err)
			c.Abort()
			return
		}
		c.Set(common.UserIDKey, payload.UserID)
		c.Set(common.UserNameKey, payload.UserName)

		c.Next()
	}
}
