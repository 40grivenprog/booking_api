package api

import (
	"context"

	"github.com/gin-gonic/gin"
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
	professionalsService "github.com/vention/booking_api/internal/services/professionals"
)

func Register(ctx context.Context, cfg *config.Config, router *gin.RouterGroup, queries *db.Queries) error {
	// Register clients API
	if err := clientsAPI.ClientsRegister(clientsAPI.ClientsHandlerParams{
		Router:         router,
		ClientsService: clientsService.NewService(queries),
	}); err != nil {
		return err
	}

	// Register professionals API
	if err := professionalsAPI.ProfessionalsRegister(professionalsAPI.ProfessionalsHandlerParams{
		Router:               router,
		ProfessionalsService: professionalsService.NewService(queries),
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
		Router:    router,
		UsersRepo: queries,
	}); err != nil {
		return err
	}

	return nil
}
