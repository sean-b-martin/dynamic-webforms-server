package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
)

type UserController struct {
}

func NewUserController(router fiber.Router) *UserController {
	controller := UserController{}
	router.Use(middleware.AllowedContentTypeWithJSON())
	router.Post("/", controller.RegisterUser)
	return &controller
}

type registerUserData struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8"`
}

func (u *UserController) RegisterUser(ctx *fiber.Ctx) error {
	var user registerUserData

	if err := parseAndValidateRequestData(ctx, nil, &user); err != nil {
		return err
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}
