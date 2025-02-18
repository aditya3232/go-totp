package models

type CreateUserRequest struct {
	Name     string `json:"name" `
	Email    string `json:"email" `
	Password string `json:"password" `
	TOTPKey  string `json:"totp_key" `
}

type UserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPKey  string `json:"totp_key"`
}
