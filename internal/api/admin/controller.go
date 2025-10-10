package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/admin"
)

// CreateProfessional handles POST /api/admin/professionals
func (h *AdminsHandler) CreateProfessional(c *gin.Context) {
	req, ok := common.BindAndValidate[CreateProfessionalRequest](c)
	if !ok {
		return
	}

	professional, err := h.adminService.CreateProfessional(c.Request.Context(), admin.CreateProfessionalInput{
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		// Check if it's a unique constraint violation
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, common.ErrorTypeConflict, common.ErrorMsgUsernameAlreadyExists, nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateProfessional, err)
		return
	}

	response := mapProfessionalToCreateProfessionalResponse(professional)
	c.JSON(http.StatusCreated, response)
}
