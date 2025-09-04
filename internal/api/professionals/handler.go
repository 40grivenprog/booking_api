package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/config"
)

type ProfessionalsHandler struct {
	cfg              *config.Config
	professionalsRepo ProfessionalsRepository
}

type ProfessionalsHandlerParams struct {
	Router            *gin.Engine
	Cfg               *config.Config
	ProfessionalsRepo ProfessionalsRepository
}

func ProfessionalsRegister(p ProfessionalsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ProfessionalsRepo == nil {
		return errors.New("missing professionals repository")
	}

	h := &ProfessionalsHandler{
		cfg:              p.Cfg,
		professionalsRepo: p.ProfessionalsRepo,
	}

	api := p.Router.Group("/api")
	{
		professionals := api.Group("/professionals")
		{
			professionals.GET("", h.GetProfessionals)
			professionals.POST("/sign_in", h.SignInProfessional)
		}
	}

	return nil
}
