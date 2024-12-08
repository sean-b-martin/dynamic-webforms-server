package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

type SchemaController struct {
	service service.SchemaService
}

func NewSchemaController(router fiber.Router, authMiddleware *middleware.JWTAuth, service service.SchemaService) *SchemaController {
	controller := SchemaController{service: service}
	router.Get("/", controller.GetFormSchemas)
	router.Get("/:schemaID", controller.GetSchema)
	router.Post("/", authMiddleware.Handle(), controller.CreateSchema)
	router.Patch("/:schemaID", authMiddleware.Handle(), controller.UpdateSchema)
	router.Delete("/:schemaID", authMiddleware.Handle(), controller.DeleteSchema)

	return &controller
}

func (s *SchemaController) GetFormSchemas(ctx *fiber.Ctx) error {
	var formID requestPathFormID
	if !parseAndValidateRequestData(ctx, &formID, nil) {
		return nil
	}

	schemas, err := s.service.GetSchemas(formID.FormID)
	if err != nil {
		return serviceErrToResponse(ctx, err)
	}

	if len(schemas) == 0 {
		return ctx.Status(fiber.StatusOK).JSON([]struct{}{})
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas)
}

func (s *SchemaController) GetSchema(ctx *fiber.Ctx) error {
	var ids requestPathFormAndSchemaID
	if !parseAndValidateRequestData(ctx, &ids, nil) {
		return nil
	}

	schema, err := s.service.GetSchema(ids.FormID, ids.SchemaID)
	if err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schema)
}

func (s *SchemaController) CreateSchema(ctx *fiber.Ctx) error {
	var formID requestPathFormID
	var schemaData requestDataCreateSchema
	if !parseAndValidateRequestData(ctx, &formID, schemaData) {
		return nil
	}

	err := s.service.CreateSchema(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), formID.FormID, model.FormSchemaModel{
		Title:    schemaData.Title,
		Version:  schemaData.Version,
		Schema:   schemaData.Schema,
		ReadOnly: schemaData.ReadOnly,
	})

	if err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (s *SchemaController) UpdateSchema(ctx *fiber.Ctx) error {
	return fiber.ErrInternalServerError
}

func (s *SchemaController) DeleteSchema(ctx *fiber.Ctx) error {
	var ids requestPathFormAndSchemaID

	if !parseAndValidateRequestData(ctx, &ids, nil) {
		return nil
	}

	if err := s.service.DeleteSchema(ctx.Locals(middleware.UserIDLocal).(uuid.UUID), ids.FormID, ids.SchemaID); err != nil {
		return serviceErrToResponse(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
