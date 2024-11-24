package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
)

type UserController struct {
}

func NewUserController(router fiber.Router, authService *middleware.JWTAuth) *UserController {
	controller := UserController{}
	router.Get("/", authService.Handle(), controller.GetCurrentLogin)

	router.Use(middleware.AllowedContentTypeWithJSON())
	router.Post("/", controller.RegisterUser)

	return &controller
}

type registerUserData struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (u *UserController) RegisterUser(ctx *fiber.Ctx) error {
	var user registerUserData

	if err := parseAndValidateRequestData(ctx, nil, &user); err != nil {
		return err
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}

func (u *UserController) GetCurrentLogin(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"userID": ctx.Locals(middleware.UserIDLocal).(string)})
}
