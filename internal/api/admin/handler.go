package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type AdminsHandler struct {
	adminsRepo AdminRepository
}

type AdminsHandlerParams struct {
	Router     *gin.Engine
	AdminsRepo AdminRepository
}

func AdminsRegister(p AdminsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.AdminsRepo == nil {
		return errors.New("missing admin repository")
	}

	h := &AdminsHandler{
		adminsRepo: p.AdminsRepo,
	}

	api := p.Router.Group("/api")
	{
		admin := api.Group("/admins")
		{
			admin.POST("/professionals", h.CreateProfessional)
		}
	}

	return nil
}
