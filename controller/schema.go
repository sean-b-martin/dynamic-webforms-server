package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

type SchemaController struct {
	service service.SchemaService
}

func NewSchemaController(router fiber.Router, authMiddleware *middleware.JWTAuth, service service.SchemaService) *SchemaController {
	controller := SchemaController{service: service}
	router.Get("/", controller.GetFormSchemas)
	router.Get("/:schemaID", controller.GetSchema)
	router.Post("/", controller.CreateSchema)
	router.Patch("/:schemaID", controller.UpdateSchema)
	router.Delete("/:schemaID", controller.DeleteSchema)

	return &controller
}

func (s *SchemaController) GetFormSchemas(ctx *fiber.Ctx) error {

	return fiber.ErrInternalServerError
}

func (s *SchemaController) GetSchema(ctx *fiber.Ctx) error {
	return fiber.ErrInternalServerError
}

func (s *SchemaController) CreateSchema(ctx *fiber.Ctx) error {
	return fiber.ErrInternalServerError
}

func (s *SchemaController) UpdateSchema(ctx *fiber.Ctx) error {
	return fiber.ErrInternalServerError
}

func (s *SchemaController) DeleteSchema(ctx *fiber.Ctx) error {
	return fiber.ErrInternalServerError
}
