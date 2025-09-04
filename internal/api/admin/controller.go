package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/config"
)

type AdminsHandler struct {
	cfg        *config.Config
	adminsRepo AdminsRepository
}

type AdminsHandlerParams struct {
	Router     *gin.Engine
	Cfg        *config.Config
	AdminsRepo AdminsRepository
}

func AdminsRegister(p AdminsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.AdminsRepo == nil {
		return errors.New("missing clients repository")
	}

	h := &AdminsHandler{
		cfg:        p.Cfg,
		adminsRepo: p.AdminsRepo,
	}

	api := p.Router.Group("/api")
	{
		clients := api.Group("/admins")
		{
			clients.POST("/professionals", h.CreateProfessional)
		}
	}

	return nil
}
