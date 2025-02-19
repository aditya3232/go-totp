package main

import (
	"fmt"
	"go-otp/app/configs"
)

func main() {
	viperConfig := configs.NewViper()
	logrus := configs.NewLogger(viperConfig)
	mysqlDB := configs.NewDatabase(viperConfig, logrus)
	redis := configs.NewRedis(viperConfig, logrus)
	validate := configs.NewValidator(viperConfig)
	app := configs.NewFiber(viperConfig)

	configs.Bootstrap(&configs.BootstrapConfig{
		MysqlDB:  mysqlDB,
		App:      app,
		Logrus:   logrus,
		Validate: validate,
		Config:   viperConfig,
		Redis:    redis,
	})

	webPort := viperConfig.GetInt("WEB_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
