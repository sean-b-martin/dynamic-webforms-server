package controller

import (
	"github.com/gofiber/fiber/v2"
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
