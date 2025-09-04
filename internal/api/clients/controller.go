package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// RegisterClient handles POST /api/clients/register
func (h *ClientsHandler) RegisterClient(c *gin.Context) {
	var req ClientRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Convert phone number to sql.NullString
	var phoneNumber sql.NullString
	if req.PhoneNumber != nil {
		phoneNumber = sql.NullString{String: *req.PhoneNumber, Valid: true}
	}

	// Convert chat_id to sql.NullInt64
	var chatID sql.NullInt64
	chatID = sql.NullInt64{Int64: req.ChatID, Valid: true}

	// Create new client
	user, err := h.clientsRepo.CreateClient(c.Request.Context(), &db.CreateClientParams{
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: phoneNumber,
		ChatID:      chatID,
	})
	if err != nil {
		// Check if it's a unique constraint violation
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, "username_taken", "The provided username is already taken", nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to create client", err)
		return
	}

	// Convert to response format
	responseUser := User{
		ID:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle optional fields
	if user.ChatID.Valid {
		responseUser.ChatID = &user.ChatID.Int64
	}
	if user.PhoneNumber.Valid {
		responseUser.PhoneNumber = &user.PhoneNumber.String
	}

	response := ClientRegisterResponse{
		User: responseUser,
	}

	c.JSON(http.StatusCreated, response)
}
