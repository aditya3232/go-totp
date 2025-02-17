package configs

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-otp/app/controllers"
	"go-otp/app/routes"
	"go-otp/app/services"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	MysqlDB  *gorm.DB
	App      *fiber.App
	Logrus   *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	//	setup services
	PingDbService := services.NewPingDbService(config.MysqlDB, config.Logrus)

	//	setup controllers
	PingDbController := controllers.NewPingDbController(PingDbService, config.Logrus)

	routeConfig := routes.RouteConfig{
		App:              config.App,
		PingDbController: PingDbController,
	}

	routeConfig.Setup()

}
