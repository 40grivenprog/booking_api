package admin

import (
	"context"

	db "github.com/vention/booking_api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the business logic operations for admin
type Service interface {
	CreateProfessional(ctx context.Context, input CreateProfessionalInput) (*db.Professional, error)
}

type service struct {
	repo AdminsRepository
}

// NewService creates a new admin service
func NewService(repo AdminsRepository) Service {
	return &service{
		repo: repo,
	}
}

// CreateProfessional creates a new professional with business logic validation
func (s *service) CreateProfessional(ctx context.Context, input CreateProfessionalInput) (*db.Professional, error) {
	// Hash password (business logic)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Prepare database parameters
	params := &db.CreateProfessionalParams{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	// Set optional phone number
	if input.PhoneNumber != "" {
		params.PhoneNumber.String = input.PhoneNumber
		params.PhoneNumber.Valid = true
	}

	// Set password hash
	params.PasswordHash.String = string(hashedPassword)
	params.PasswordHash.Valid = true

	// Create professional in database
	professional, err := s.repo.CreateProfessional(ctx, params)
	if err != nil {
		return nil, err
	}

	return professional, nil
}
