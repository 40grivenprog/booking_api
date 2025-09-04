package api

import (
	"context"

	"github.com/gin-gonic/gin"
	adminAPI "github.com/vention/booking_api/internal/api/admin"
	clientsAPI "github.com/vention/booking_api/internal/api/clients"
	professionalsAPI "github.com/vention/booking_api/internal/api/professionals"
	"github.com/vention/booking_api/internal/config"
	db "github.com/vention/booking_api/internal/repository"
)

func Register(ctx context.Context, cfg *config.Config, router *gin.Engine, queries *db.Queries) error {
	// Register clients API
	if err := clientsAPI.ClientsRegister(clientsAPI.ClientsHandlerParams{
		Router:      router,
		Cfg:         cfg,
		ClientsRepo: queries,
	}); err != nil {
		return err
	}

	// Register professionals API
	if err := professionalsAPI.ProfessionalsRegister(professionalsAPI.ProfessionalsHandlerParams{
		Router:            router,
		Cfg:               cfg,
		ProfessionalsRepo: queries,
	}); err != nil {
		return err
	}

	// Register admin API
	if err := adminAPI.AdminsRegister(adminAPI.AdminsHandlerParams{
		Router:     router,
		Cfg:        cfg,
		AdminsRepo: queries,
	}); err != nil {
		return err
	}

	return nil
}
