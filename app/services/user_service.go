package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-otp/app/entities"
	"go-otp/app/interfaces"
	"go-otp/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Mysql    *gorm.DB
	Logrus   *logrus.Logger
	Validate *validator.Validate
	UserRepo interfaces.IUserRepository
}

func NewUserService(mysql *gorm.DB, logrus *logrus.Logger, validate *validator.Validate, userRepo interfaces.IUserRepository) *UserService {
	return &UserService{
		Mysql:    mysql,
		Logrus:   logrus,
		Validate: validate,
		UserRepo: userRepo,
	}
}

func (s *UserService) Get(ctx context.Context) (users []models.UserResponse, err error) {
	tx := s.Mysql.WithContext(ctx)

	userEntities, err := s.UserRepo.Get(tx)
	if err != nil {
		s.Logrus.WithError(err).Error("error getting users")
		return nil, fiber.ErrInternalServerError
	}

	users = make([]models.UserResponse, len(userEntities))
	for i, user := range userEntities {
		users[i] = models.UserResponse{
			Name:    user.Name,
			Email:   user.Email,
			TOTPKey: user.TOTPKey,
		}
	}

	return users, nil
}

func (s *UserService) FindByEmail(ctx context.Context, request *models.FindUserByEmailRequest) (users *models.UserResponse, err error) {
	tx := s.Mysql.WithContext(ctx)

	if err := s.Validate.Struct(request); err != nil {
		s.Logrus.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	newReq := &models.FindUserByEmailRequest{
		Email: request.Email,
	}

	usersEntities, err := s.UserRepo.FindByEmail(tx, newReq)
	if err != nil {
		s.Logrus.WithError(err).Error("error getting users")
		return nil, fiber.ErrInternalServerError
	}

	users = &models.UserResponse{
		Name:    usersEntities.Name,
		Email:   usersEntities.Email,
		TOTPKey: usersEntities.TOTPKey,
	}

	return users, nil
}

func (s *UserService) Create(ctx context.Context, request *models.CreateUserRequest) (err error) {
	tx := s.Mysql.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Logrus.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		s.Logrus.WithError(err).Error("error hashing password")
		return fiber.ErrInternalServerError
	}

	user := &entities.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		TOTPKey:  request.TOTPKey,
	}

	if err := s.UserRepo.Create(tx, user); err != nil {
		s.Logrus.WithError(err).Error("error creating user")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Logrus.WithError(err).Error("error creating user")
		return fiber.ErrInternalServerError
	}

	return nil
}
