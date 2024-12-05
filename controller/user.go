package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

type UserController struct {
	service service.UserService
}

func NewUserController(router fiber.Router, authService *middleware.JWTAuth, userService service.UserService) *UserController {
	controller := UserController{service: userService}
	router.Get("/login", authService.Handle(), controller.GetCurrentLogin)

	router.Use(middleware.AllowedContentTypeWithJSON())
	router.Post("/register", controller.RegisterUser)
	router.Post("/login", controller.LoginUser)

	return &controller
}

type userRequestData struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (u *UserController) GetCurrentLogin(ctx *fiber.Ctx) error {
	user, err := u.service.GetUserById(ctx.Locals(middleware.UserIDLocal).(uuid.UUID))
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "an error has occurred")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": user.ID, "username": user.Username})
}

func (u *UserController) LoginUser(ctx *fiber.Ctx) error {
	var user userRequestData
	if err, ok := parseAndValidateRequestData(ctx, nil, &user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if !ok {
		return nil
	}

	token, err := u.service.LoginUser(model.UserModel{Username: user.Username, Password: user.Password})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("invalid username or password")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (u *UserController) RegisterUser(ctx *fiber.Ctx) error {
	var user userRequestData
	if err, ok := parseAndValidateRequestData(ctx, nil, &user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if !ok {
		return nil
	}

	if err := u.service.RegisterUser(model.UserModel{
		Username: user.Username,
		Password: user.Password,
	}); err != nil {
		return err
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}
