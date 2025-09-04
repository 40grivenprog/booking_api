package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// GetProfessionals handles GET /api/professionals
func (h *ProfessionalsHandler) GetProfessionals(c *gin.Context) {
	professionals, err := h.professionalsRepo.GetProfessionals(c.Request.Context())
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to retrieve professionals", err)
		return
	}

	// Convert to response format
	var responseUsers []User
	for _, prof := range professionals {
		user := User{
			ID:        prof.ID.String(),
			Username:  prof.Username,
			FirstName: prof.FirstName,
			LastName:  prof.LastName,
			UserType:  string(prof.UserType),
			CreatedAt: prof.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: prof.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Handle optional fields
		if prof.ChatID.Valid {
			user.ChatID = &prof.ChatID.Int64
		}
		if prof.PhoneNumber.Valid {
			user.PhoneNumber = &prof.PhoneNumber.String
		}

		responseUsers = append(responseUsers, user)
	}

	response := GetProfessionalsResponse{
		Professionals: responseUsers,
	}

	c.JSON(http.StatusOK, response)
}

// SignInProfessional handles POST /api/professionals/sign_in
func (h *ProfessionalsHandler) SignInProfessional(c *gin.Context) {
	var req ProfessionalSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Get user by username
	user, err := h.professionalsRepo.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Check if user is a professional
	if user.UserType != "professional" {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "User is not a professional", nil)
		return
	}

	// Check if password hash exists
	if !user.PasswordHash.Valid {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Update user's chat_id after successful authentication
	updatedUser, err := h.professionalsRepo.UpdateUserChatID(c.Request.Context(), &db.UpdateUserChatIDParams{
		ID:     user.ID,
		ChatID: sql.NullInt64{Int64: req.ChatID, Valid: true},
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to update user chat_id", err)
		return
	}

	// Convert to response format
	responseUser := User{
		ID:        updatedUser.ID.String(),
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		UserType:  string(updatedUser.UserType),
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle optional fields
	if updatedUser.ChatID.Valid {
		responseUser.ChatID = &updatedUser.ChatID.Int64
	}
	if updatedUser.PhoneNumber.Valid {
		responseUser.PhoneNumber = &updatedUser.PhoneNumber.String
	}

	response := ProfessionalSignInResponse{
		User:    responseUser,
	}

	c.JSON(http.StatusOK, response)
}
