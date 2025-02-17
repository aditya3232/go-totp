package services

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PingDbService struct {
	Mysql  *gorm.DB
	Logrus *logrus.Logger
}

func NewPingDbService(mysql *gorm.DB, logrus *logrus.Logger) *PingDbService {
	return &PingDbService{
		Mysql:  mysql,
		Logrus: logrus,
	}
}

func (s *PingDbService) Ping(ctx context.Context) error {
	mysqlDB, err := s.Mysql.DB()
	if err != nil {
		s.Logrus.WithError(err).Error("Failed to get DB instance")
		return fiber.ErrInternalServerError
	}

	if err := mysqlDB.PingContext(ctx); err != nil {
		s.Logrus.WithError(err).Error("Database ping failed")
		return fiber.ErrInternalServerError
	}

	s.Logrus.Info("Database mysql ping successful")
	return nil
}
