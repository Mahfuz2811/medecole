package models

// LoginRequest represents the login request payload
type LoginRequest struct {
	MSISDN   string `json:"msisdn" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	MSISDN   string `json:"msisdn" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
