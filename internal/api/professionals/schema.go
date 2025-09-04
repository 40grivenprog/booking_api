package api

// ProfessionalSignInRequest represents the request body for professional sign in
type ProfessionalSignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ChatID   int64  `json:"chat_id" binding:"required"`
}

// ProfessionalSignInResponse represents the response for professional sign in
type ProfessionalSignInResponse struct {
	User    User   `json:"user"`
}

// GetProfessionalsResponse represents the response for getting all professionals
type GetProfessionalsResponse struct {
	Professionals []User `json:"professionals"`
}

// User represents a user in API responses (using SQLC generated model)
type User struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	UserType    string  `json:"user_type"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
