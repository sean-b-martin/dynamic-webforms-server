package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

type FormController struct {
	service service.FormService
}

func NewFormController(router fiber.Router, authMiddleware *middleware.JWTAuth, service service.FormService) *FormController {
	controller := FormController{service: service}
	router.Get("/my-forms", authMiddleware.Handle(), controller.GetMyForms)
	router.Get("/:formID", controller.GetForm)
	router.Get("/", controller.GetForms)
	router.Post("/", authMiddleware.Handle(), controller.CreateForm)
	router.Patch("/:formID", authMiddleware.Handle(), controller.UpdateForm)
	router.Delete("/:formID", authMiddleware.Handle(), controller.DeleteForm)
	return &controller
}

type formIDPath struct {
	FormID uuid.UUID `json:"formID" validate:"required,uuid"`
}

type FormRequest struct {
	Title string `json:"title" validate:"required,min=1,max=256"`
}

func (c *FormController) GetForms(ctx *fiber.Ctx) error {
	forms, err := c.service.GetForms()

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if forms == nil {
		return ctx.Status(fiber.StatusOK).JSON([]struct{}{})
	}

	return ctx.Status(fiber.StatusOK).JSON(forms)
}

func (c *FormController) GetForm(ctx *fiber.Ctx) error {
	var formID formIDPath
	if ok := parseAndValidateRequestData(ctx, &formID, nil); !ok {
		return nil
	}

	form, err := c.service.GetForm(formID.FormID)
	if err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.JSON(form)
}

func (c *FormController) GetMyForms(ctx *fiber.Ctx) error {
	forms, err := c.service.GetFormsOfUser(ctx.Locals(middleware.UserIDLocal).(uuid.UUID))
	if err != nil {
		return serviceErrToResponse(ctx, err)
	}

	if forms == nil {
		return ctx.Status(fiber.StatusOK).JSON([]struct{}{})
	}

	return ctx.Status(fiber.StatusOK).JSON(forms)
}

func (c *FormController) CreateForm(ctx *fiber.Ctx) error {
	var form FormRequest
	if ok := parseAndValidateRequestData(ctx, nil, &form); !ok {
		return nil
	}

	if err := c.service.CreateForm(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), form.Title); err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *FormController) UpdateForm(ctx *fiber.Ctx) error {
	var formID formIDPath
	var form FormRequest

	if ok := parseAndValidateRequestData(ctx, &formID, &form); !ok {
		return nil
	}

	if err := c.service.UpdateForm(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), formID.FormID, form.Title); err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *FormController) DeleteForm(ctx *fiber.Ctx) error {
	var formID formIDPath
	if ok := parseAndValidateRequestData(ctx, nil, &formID); !ok {
		return nil
	}

	if err := c.service.DeleteForm(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), formID.FormID); err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
