package interfaces

import (
	"context"
	"go-otp/app/entities"
	"go-otp/app/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Get(mysql *gorm.DB) (users []entities.User, err error)
	FindByEmail(mysql *gorm.DB, email string) (user entities.User, err error)
	UpdateColumn(mysql *gorm.DB, id int, column string, value interface{}) (err error)
	Create(mysql *gorm.DB, user entities.User) (err error)
}

type IUserService interface {
	Get(ctx context.Context) (users []models.UserResponse, err error)
	Create(ctx context.Context, request *models.CreateUserRequest) (err error)
}
