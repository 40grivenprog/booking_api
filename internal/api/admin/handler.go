package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/admin"
)

// AdminsHandler handles HTTP requests for admin operations
type AdminsHandler struct {
	adminService admin.Service
}

// NewAdminsHandler creates a new handler with dependency injection
func NewAdminsHandler(service admin.Service) *AdminsHandler {
	return &AdminsHandler{
		adminService: service,
	}
}

// AdminsHandlerParams defines the parameters for the AdminsHandler
type AdminsHandlerParams struct {
	Router       *gin.RouterGroup
	AdminService admin.Service
}

// AdminsRegister registers the AdminsHandler with the router
func AdminsRegister(p AdminsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.AdminService == nil {
		return errors.New("missing admin service")
	}

	h := NewAdminsHandler(p.AdminService)

	admin := p.Router.Group("/admins")
	{
		admin.POST("/professionals", h.CreateProfessional)
	}

	return nil
}
