package models

// User represents a user in the system
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Never serialize password hash to JSON
	Role         string `json:"role"`
}

// UserForm represents the data needed to create/update a user
type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

