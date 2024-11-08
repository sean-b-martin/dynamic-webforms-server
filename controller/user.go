package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/sean-b-martin/dynamic-webforms-server/validation"
)

type UserController struct {
}

func NewUserController(app fiber.Router) *UserController {
	controller := UserController{}
	app.Post("/", controller.RegisterUser)
	return &controller
}

type registerUserData struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8"`
}

func (u *UserController) RegisterUser(ctx *fiber.Ctx) error {
	var user registerUserData

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := validation.Validate(user); len(err) > 0 {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": err})
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}
