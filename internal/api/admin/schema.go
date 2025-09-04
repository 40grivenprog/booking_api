package api

// CreateProfessionalRequest represents the request to create a professional
type CreateProfessionalRequest struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// CreateProfessionalResponse represents the response after creating a professional
type CreateProfessionalResponse struct {
	User    User   `json:"user"`
}

// User represents a user in the response
type User struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	UserType    string  `json:"user_type"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
