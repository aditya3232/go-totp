package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-otp/app/interfaces"
	"go-otp/app/objects"
)

type PingDbController struct {
	PingDbService interfaces.IPingDbService
	Logrus        *logrus.Logger
}

func NewPingDbController(pingDbService interfaces.IPingDbService, logrus *logrus.Logger) *PingDbController {
	return &PingDbController{
		PingDbService: pingDbService,
		Logrus:        logrus,
	}
}

func (c *PingDbController) Ping(ctx *fiber.Ctx) error {
	if err := c.PingDbService.Ping(ctx.UserContext()); err != nil {
		c.Logrus.WithError(err).Error("error ping mysql")
		return ctx.Status(fiber.StatusInternalServerError).JSON(objects.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(objects.Response{
		Message: "Database mysql ping successful",
		Data:    nil,
	})
}
