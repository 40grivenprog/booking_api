package token

// Maker is an interface for managing tokens
type Maker interface {
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
