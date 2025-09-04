package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/config"
)

type ClientsHandler struct {
	cfg        *config.Config
	clientsRepo ClientsRepository
}

type ClientsHandlerParams struct {
	Router      *gin.Engine
	Cfg         *config.Config
	ClientsRepo ClientsRepository
}

func ClientsRegister(p ClientsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ClientsRepo == nil {
		return errors.New("missing clients repository")
	}

	h := &ClientsHandler{
		cfg:        p.Cfg,
		clientsRepo: p.ClientsRepo,
	}

	api := p.Router.Group("/api")
	{
		clients := api.Group("/clients")
		{
			clients.POST("/register", h.RegisterClient)
		}
	}

	return nil
}
