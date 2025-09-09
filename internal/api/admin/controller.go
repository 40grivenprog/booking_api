package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateProfessional handles POST /api/admin/professionals
func (h *AdminsHandler) CreateProfessional(c *gin.Context) {
	var req CreateProfessionalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidRequestBody, err)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeInternal, common.ErrorMsgInternalServerError, err)
		return
	}

	// Convert phone number to sql.NullString
	var phoneNumber sql.NullString
	if req.PhoneNumber != "" {
		phoneNumber = sql.NullString{String: req.PhoneNumber, Valid: true}
	}

	// Convert chat_id to sql.NullInt64 (can be NULL for admin-created professionals)
	var chatID sql.NullInt64
	// chatID will be NULL by default

	// Create new professional
	user, err := h.adminsRepo.CreateProfessional(c, &db.CreateProfessionalParams{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  phoneNumber,
		PasswordHash: sql.NullString{String: string(hashedPassword), Valid: true},
		ChatID:       chatID,
	})
	if err != nil {
		// Check if it's a unique constraint violation
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, common.ErrorTypeConflict, common.ErrorMsgUsernameAlreadyExists, nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateProfessional, err)
		return
	}

	// Convert to response format
	responseUser := User{
		ID:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserType:  "professional",
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle optional fields
	if user.PhoneNumber.Valid {
		responseUser.PhoneNumber = &user.PhoneNumber.String
	}

	response := CreateProfessionalResponse{
		User: responseUser,
	}

	c.JSON(http.StatusCreated, response)
}
