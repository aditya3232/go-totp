package repositories

import (
	"github.com/sirupsen/logrus"
	"go-otp/app/entities"
	"go-otp/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	Logrus *logrus.Logger
}

func NewUserRepository(logrus *logrus.Logger) *UserRepository {
	return &UserRepository{
		Logrus: logrus,
	}
}

func (r *UserRepository) Get(mysql *gorm.DB) (users []entities.User, err error) {
	err = mysql.Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByEmail(mysql *gorm.DB, request *models.FindUserByEmailRequest) (user *entities.User, err error) {
	err = mysql.Where("email = ?", request.Email).First(&user).Error
	return
}

func (r *UserRepository) UpdateColumn(mysql *gorm.DB, id int, column string, value interface{}) (err error) {
	err = mysql.Model(&entities.User{}).Where("id = ?", id).Update(column, value).Error
	return
}

func (r *UserRepository) Create(mysql *gorm.DB, request *entities.User) (err error) {
	err = mysql.Create(&request).Error
	return
}
