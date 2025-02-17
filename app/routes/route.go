package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-otp/app/controllers"
)

type RouteConfig struct {
	App              *fiber.App
	Log              *logrus.Logger
	PingDbController *controllers.PingDbController
}

func (c *RouteConfig) Setup() {
	c.App.Use(c.recoverPanic)
	c.SetupGuestRoute()
	//c.SetupAuthRoute()
}

func (c *RouteConfig) recoverPanic(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic occurred: %v", r)
			c.Log.WithError(err).Error("Panic occured")
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}()

	return ctx.Next()
}

func (rc *RouteConfig) SetupGuestRoute() {
	GuestGroup := rc.App.Group("/api")

	GuestGroup.Get("/ping", rc.PingDbController.Ping)

}
