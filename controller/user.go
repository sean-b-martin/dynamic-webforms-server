package controller

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

type UserController struct {
	service service.UserService
}

func NewUserController(router fiber.Router, authMiddleware *middleware.JWTAuth, userService service.UserService) *UserController {
	controller := UserController{service: userService}
	router.Get("/login", authMiddleware.Handle(), controller.GetCurrentLogin)
	router.Delete("/", authMiddleware.Handle(), controller.DeleteUser)

	router.Use(middleware.AllowedContentTypeWithJSON())
	router.Post("/register", controller.RegisterUser)
	router.Post("/login", controller.LoginUser)
	router.Patch("/", authMiddleware.Handle(), controller.UpdateUser)

	return &controller
}

type usernameRequestData struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}

type passwordRequestData struct {
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type userRequestData struct {
	usernameRequestData
	passwordRequestData
}

type userUpdateData struct {
	passwordRequestData
}

func (u *UserController) GetCurrentLogin(ctx *fiber.Ctx) error {
	user, err := u.service.GetUserById(ctx.Locals(middleware.UserIDLocal).(uuid.UUID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user was deleted"})
		}

		return fiber.NewError(fiber.StatusInternalServerError, "an error has occurred")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": user.ID, "username": user.Username})
}

func (u *UserController) LoginUser(ctx *fiber.Ctx) error {
	var user userRequestData
	if ok := parseAndValidateRequestData(ctx, nil, &user); !ok {
		return nil
	}

	token, err := u.service.LoginUser(model.UserModel{Username: user.Username, Password: user.Password})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (u *UserController) RegisterUser(ctx *fiber.Ctx) error {
	var user userRequestData
	if ok := parseAndValidateRequestData(ctx, nil, &user); !ok {
		return nil
	}

	if err := u.service.RegisterUser(model.UserModel{
		Username: user.Username,
		Password: user.Password,
	}); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var user userUpdateData
	if ok := parseAndValidateRequestData(ctx, nil, &user); !ok {
		return nil
	}

	if err := u.service.UpdateUser(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (u *UserController) DeleteUser(ctx *fiber.Ctx) error {
	if err := u.service.DeleteUser(ctx.Locals(middleware.UserIDLocal).(uuid.UUID)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.ErrNotFound
		}

		return fiber.ErrBadRequest
	}

	return ctx.SendStatus(fiber.StatusOK)
}
