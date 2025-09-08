package api

// GetUserByChatIDResponse represents the response for getting a user by chat_id
type GetUserByChatIDResponse struct {
	User User `json:"user"`
}

// User represents a user in API responses (unified for both clients and professionals)
type User struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Role        string  `json:"role"` // "client" or "professional"
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
