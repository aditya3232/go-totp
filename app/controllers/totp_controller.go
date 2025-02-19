package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-otp/app/interfaces"
	"go-otp/app/models"
	"net/http"
)

type TotpController struct {
	Logrus      *logrus.Logger
	TotpService interfaces.ITotpService
}

func NewTotpControler(logrus *logrus.Logger, totpService interfaces.ITotpService) *TotpController {
	return &TotpController{
		Logrus:      logrus,
		TotpService: totpService,
	}
}

func (c *TotpController) InitiateEnrollment(ctx *fiber.Ctx) error {
	request := new(models.TotpRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Logrus.WithError(err).Error("error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err))
	}

	result, err := c.TotpService.InitiateEnrollment(ctx.UserContext(), request)
	if err != nil {
		c.Logrus.WithError(err).Error("error initiate enrollment")
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Message: http.StatusText(fiber.StatusOK),
		Data:    result,
	})
}

func (c *TotpController) ConfirmEnrollment(ctx *fiber.Ctx) error {
	request := new(models.TotpRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Logrus.WithError(err).Error("error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err))
	}

	if err := c.TotpService.ConfirmEnrollment(ctx.UserContext(), request); err != nil {
		c.Logrus.WithError(err).Error("error confirm enrollment")
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Message: http.StatusText(fiber.StatusOK),
	})
}

func (c *TotpController) VerifyTotp(ctx *fiber.Ctx) error {
	request := new(models.TotpRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Logrus.WithError(err).Error("error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err))
	}

	if err := c.TotpService.VerifyTotp(ctx.UserContext(), request); err != nil {
		c.Logrus.WithError(err).Error("error verify totp")
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Message: http.StatusText(fiber.StatusOK),
	})

}
