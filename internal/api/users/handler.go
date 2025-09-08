package api

import (
	"github.com/gin-gonic/gin"
)

// UsersHandlerParams contains the parameters needed to register users handlers
type UsersHandlerParams struct {
	Router    *gin.Engine
	UsersRepo UsersRepository
}

// UsersRegister registers all users-related routes
func UsersRegister(params UsersHandlerParams) error {
	// Create controller
	controller := NewUsersController(params.UsersRepo)

	// Create API group
	api := params.Router.Group("/api")
	users := api.Group("/users")

	// Register routes
	users.GET("/:chat_id", controller.GetUserByChatID)

	return nil
}
