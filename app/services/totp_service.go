package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-otp/app/interfaces"
	"go-otp/app/models"
	"gorm.io/gorm"
	"image/png"
	"time"
)

type TotpService struct {
	Mysql    *gorm.DB
	Redis    *redis.Client
	Logrus   *logrus.Logger
	Validate *validator.Validate
	UserRepo interfaces.IUserRepository
}

func NewTotpService(mysql *gorm.DB, redis *redis.Client, logrus *logrus.Logger, validate *validator.Validate, userRepo interfaces.IUserRepository) *TotpService {
	return &TotpService{
		Mysql:    mysql,
		Redis:    redis,
		Logrus:   logrus,
		Validate: validate,
		UserRepo: userRepo,
	}
}

func (s *TotpService) InitiateEnrollment(ctx context.Context, request *models.TotpRequest) (response *models.TotpEnrollmentResponse, err error) {
	tx := s.Mysql.WithContext(ctx)

	newReq := &models.FindUserByEmailRequest{
		Email: request.Email,
	}

	_, err = s.UserRepo.FindByEmail(tx, newReq)
	if err != nil {
		s.Logrus.WithError(err).Error("error find email")
		return nil, fiber.ErrNotFound
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "aditya3232",
		AccountName: request.Email,
	})
	if err != nil {
		s.Logrus.WithError(err).Error("error find generate key totp")
		return nil, fiber.ErrInternalServerError
	}

	// Create QR code
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		s.Logrus.WithError(err).Error("error generate image QRCode")
		return nil, fiber.ErrInternalServerError
	}

	// convert to base64
	err = png.Encode(&buf, img)
	if err != nil {
		s.Logrus.WithError(err).Error("error encode base64 image QRCode")
		return nil, fiber.ErrInternalServerError
	}
	qrcode := base64.StdEncoding.EncodeToString(buf.Bytes())

	response = &models.TotpEnrollmentResponse{
		QRCode:    qrcode,
		SecretKey: key.Secret(),
	}

	s.Redis.Set(ctx, fmt.Sprintf("enrollment:%s", request.Email), key.Secret(), 0)

	return response, nil
}

func (s *TotpService) ConfirmEnrollment(ctx context.Context, request *models.TotpRequest) (err error) {
	tx := s.Mysql.WithContext(ctx).Begin()
	defer tx.Rollback()

	newReq := &models.FindUserByEmailRequest{
		Email: request.Email,
	}

	user, err := s.UserRepo.FindByEmail(tx, newReq)
	if err != nil {
		s.Logrus.WithError(err).Error("error find email")
		return fiber.ErrNotFound
	}

	enrollment, err := s.Redis.Get(ctx, "enrollment:"+request.Email).Result()
	if err != nil {
		s.Logrus.WithError(err).Error("key redis not found")
		return fiber.ErrInternalServerError
	}

	// validate TOTP
	valid := totp.Validate(request.PassCode, enrollment)
	if !valid {
		s.Logrus.Error("invalid passcode")
		return fiber.ErrInternalServerError
	}

	if err := s.UserRepo.UpdateColumn(tx, user.ID, "totp_key", enrollment); err != nil {
		s.Logrus.WithError(err).Error("error update user")
		return fiber.ErrInternalServerError
	}

	// del cache
	s.Redis.Del(ctx, "enrollment:"+request.Email)

	// set cooldown
	s.Redis.Set(ctx, "cooldown:"+request.Email+":"+request.PassCode, true, 1*time.Minute)

	if err := tx.Commit().Error; err != nil {
		s.Logrus.WithError(err).Error("error commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *TotpService) VerifyTotp(ctx context.Context, request *models.TotpRequest) (err error) {
	tx := s.Mysql.WithContext(ctx)

	newReq := &models.FindUserByEmailRequest{
		Email: request.Email,
	}

	user, err := s.UserRepo.FindByEmail(tx, newReq)
	if err != nil {
		s.Logrus.WithError(err).Error("error find email")
		return fiber.ErrNotFound
	}

	// check cooldown
	found, err := s.Redis.Get(ctx, "cooldown:"+request.Email+":"+request.PassCode).Result()
	if err != nil {
		s.Logrus.WithError(err).Error("key redis not found")
		return fiber.ErrInternalServerError
	}

	if found == "true" {
		s.Logrus.Error("invalid passcode")
		return fiber.ErrInternalServerError
	}

	// validate otp
	valid := totp.Validate(request.PassCode, user.TOTPKey)
	if !valid {
		s.Logrus.Error("invalid passcode")
		return fiber.ErrInternalServerError
	}

	s.Redis.Set(ctx, "cooldown:"+request.Email+":"+request.PassCode, true, 1*time.Minute)

	return nil
}
