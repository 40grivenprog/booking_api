package api

// ClientRegisterRequest represents the request body for client registration
type ClientRegisterRequest struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	ChatID      int64   `json:"chat_id" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// ClientRegisterResponse represents the response for client registration
type ClientRegisterResponse struct {
	User User `json:"user"`
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
