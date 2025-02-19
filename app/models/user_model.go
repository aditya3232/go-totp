package models

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	TOTPKey  string `json:"totp_key" validate:"required,max=100"`
}

type FindUserByEmailRequest struct {
	Email string `json:"email" validate:"required,max=100"`
}

type UserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPKey  string `json:"totp_key"`
}
