package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/validation"
)

func parseAndValidateRequestData(ctx *fiber.Ctx, paramsOut interface{}, bodyOut interface{}) bool {
	if paramsOut != nil {
		if err := ctx.ParamsParser(paramsOut); err != nil {
			_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			return false
		}

		if validationErrors := validation.Validate(paramsOut); len(validationErrors) > 0 {
			_ = ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
			return false
		}
	}

	if bodyOut != nil {
		if err := ctx.BodyParser(bodyOut); err != nil {
			_ = ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			return false
		}

		if validationErrors := validation.Validate(bodyOut); len(validationErrors) > 0 {
			_ = ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
			return false
		}
	}

	return true
}

// definitions for path and request data

type requestDataUsername struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}

type requestDataPassword struct {
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type requestDataUser struct {
	requestDataUsername
	requestDataPassword
}

type requestDataUpdateUser struct {
	requestDataPassword
}

type requestDataTitle struct {
	Title string `json:"title" validate:"required,min=1,max=256"`
}

type requestDataCreateSchema struct {
	requestDataTitle
	Version  string `json:"version" validate:"required,min=1,max=64"`
	Schema   []byte `json:"schema"`
	ReadOnly bool   `json:"readOnly"`
}

type requestDataUpdateSchema struct {
	Title    *string `json:"title,omitempty" validate:"min=1,max=256"`
	Schema   *[]byte `json:"schema,omitempty"`
	ReadOnly *bool   `json:"readOnly,omitempty"`
}

// path structs

type requestPathFormID struct {
	FormID uuid.UUID `json:"formID" validate:"required,uuid"`
}

type requestPathSchemaID struct {
	SchemaID uuid.UUID `json:"schemaID" validate:"required,uuid"`
}

type requestPathFormAndSchemaID struct {
	requestPathFormID
	requestPathSchemaID
}
