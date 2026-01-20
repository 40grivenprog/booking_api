package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	adminAPI "github.com/vention/booking_api/internal/api/admin"
	appointmentsAPI "github.com/vention/booking_api/internal/api/appointments"
	clientsAPI "github.com/vention/booking_api/internal/api/clients"
	professionalsAPI "github.com/vention/booking_api/internal/api/professionals"
	usersAPI "github.com/vention/booking_api/internal/api/users"
	"github.com/vention/booking_api/internal/config"
	db "github.com/vention/booking_api/internal/repository"
	adminService "github.com/vention/booking_api/internal/services/admin"
	appointmentsService "github.com/vention/booking_api/internal/services/appointments"
	clientsService "github.com/vention/booking_api/internal/services/clients"
	"github.com/vention/booking_api/internal/services/notifications"
	professionalsService "github.com/vention/booking_api/internal/services/professionals"
	"github.com/vention/booking_api/internal/token"
)

func Register(ctx context.Context, cfg *config.Config, router *gin.RouterGroup, queries *db.Queries, tokenMaker token.Maker) error {
	// Register clients API
	notificationsService, err := notifications.NewService(cfg.TelegramBotToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create notifications service")
		return err
	}
	if err := clientsAPI.ClientsRegister(clientsAPI.ClientsHandlerParams{
		TokenMaker:           tokenMaker,
		Router:               router,
		ClientsService:       clientsService.NewService(queries),
		NotificationsService: notificationsService,
	}); err != nil {
		return err
	}

	// Register professionals API
	if err := professionalsAPI.ProfessionalsRegister(professionalsAPI.ProfessionalsHandlerParams{
		TokenMaker:           tokenMaker,
		Router:               router,
		ProfessionalsService: professionalsService.NewService(queries),
		NotificationsService: notificationsService,
	}); err != nil {
		return err
	}

	// Register admin API
	if err := adminAPI.AdminsRegister(adminAPI.AdminsHandlerParams{
		Router:       router,
		AdminService: adminService.NewService(queries),
	}); err != nil {
		return err
	}

	// Register appointments API
	if err := appointmentsAPI.AppointmentsRegister(appointmentsAPI.AppointmentsHandlerParams{
		Router:              router,
		AppointmentsService: appointmentsService.NewService(queries),
	}); err != nil {
		return err
	}

	// Register users API
	if err := usersAPI.UsersRegister(usersAPI.UsersHandlerParams{
		Router:     router,
		UsersRepo:  queries,
		TokenMaker: tokenMaker,
	}); err != nil {
		return err
	}

	return nil
}
