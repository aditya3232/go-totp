package models

type TotpRequest struct {
	Email    string `json:"email" validate:"max=100"`
	PassCode string `json:"passcode" validate:"max=100"`
}

type TotpVerifyRequest struct {
	PassCode string `json:"passcode" validate:"required,max=100"`
}

type TotpEnrollmentResponse struct {
	QRCode    string `json:"qrcode"`
	SecretKey string `json:"secret"`
}
