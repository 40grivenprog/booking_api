package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

func mapProfessionalToCreateProfessionalResponse(professional *db.Professional) CreateProfessionalResponse {
	responseUser := User{
		ID:          professional.ID.String(),
		Username:    professional.Username,
		FirstName:   professional.FirstName,
		LastName:    professional.LastName,
		UserType:    common.UserTypeProfessional,
		PhoneNumber: common.FromNullString(professional.PhoneNumber),
		CreatedAt:   common.FormatTimeWithTimezone(professional.CreatedAt),
		UpdatedAt:   common.FormatTimeWithTimezone(professional.UpdatedAt),
	}

	return CreateProfessionalResponse{
		User: responseUser,
	}
}
