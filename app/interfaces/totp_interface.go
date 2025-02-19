package interfaces

import (
	"context"
	"go-otp/app/models"
)

type ITotpService interface {
	InitiateEnrollment(ctx context.Context, request *models.TotpRequest) (response *models.TotpEnrollmentResponse, err error)
	ConfirmEnrollment(ctx context.Context, request *models.TotpRequest) (err error)
	VerifyTotp(ctx context.Context, request *models.TotpRequest) (err error)
}
