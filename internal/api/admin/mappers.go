package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// mapProfessionalToCreateProfessionalResponse maps a professional to a CreateProfessionalResponse
func mapProfessionalToCreateProfessionalResponse(professional *db.Professional) CreateProfessionalResponse {
	responseUser := CreateProfessionalResponseItem{
		ID:        professional.ID.String(),
		Username:  professional.Username,
		FirstName: professional.FirstName,
		LastName:  professional.LastName,
		UserType:  common.UserTypeProfessional,
		CreatedAt: common.FormatTimeWithTimezone(professional.CreatedAt),
		UpdatedAt: common.FormatTimeWithTimezone(professional.UpdatedAt),
	}

	return CreateProfessionalResponse{
		CreateProfessionalResponseItem: responseUser,
	}
}
