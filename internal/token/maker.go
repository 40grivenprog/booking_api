package token

import "github.com/google/uuid"

// Maker is an interface for managing tokens
type Maker interface {
	VerifyToken(token string) (*Payload, error)
	CreateToken(userID uuid.UUID) (string, error)
}
