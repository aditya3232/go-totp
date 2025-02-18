package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-otp/app/interfaces"
	"go-otp/app/models"
)

type UserController struct {
	Logrus      *logrus.Logger
	UserService interfaces.IUserService
}

func NewUserController(logrus *logrus.Logger, userService interfaces.IUserService) *UserController {
	return &UserController{
		Logrus:      logrus,
		UserService: userService,
	}
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	request := new(models.CreateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Logrus.WithError(err).Error("error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err))
	}

	if err := c.UserService.Create(ctx.UserContext(), request); err != nil {
		c.Logrus.WithError(err).Error("error creating user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.Response{
		Message: "data user created",
		Data:    nil,
	})
}

func (c *UserController) Get(ctx *fiber.Ctx) error {
	users, err := c.UserService.Get(ctx.UserContext())
	if err != nil {
		c.Logrus.WithError(err).Error("error get users")
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Message: "success",
		Data:    users,
	})
}
