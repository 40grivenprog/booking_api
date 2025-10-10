package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindAndValidate validates JSON request body and handles errors automatically
// Returns the validated request and a boolean indicating success
// If validation fails, error response is sent automatically and false is returned
func BindAndValidate[T any](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgInvalidRequestBody, err)
		return req, false
	}
	return req, true
}
