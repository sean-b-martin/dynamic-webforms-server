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

	return &controller
}
