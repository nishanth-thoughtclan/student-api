package utils

// AuthRequest represents the payload for authentication requests
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
