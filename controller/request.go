package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/validation"
)

func parseAndValidateRequestData(ctx *fiber.Ctx, paramsOut interface{}, bodyOut interface{}) error {
	if paramsOut != nil {
		if err := ctx.ParamsParser(paramsOut); err != nil {
			return err
		}

		if validationErrors := validation.Validate(paramsOut); len(validationErrors) > 0 {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
		}
	}

	if bodyOut != nil {
		if err := ctx.BodyParser(bodyOut); err != nil {
			return err
		}

		if validationErrors := validation.Validate(bodyOut); len(validationErrors) > 0 {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
		}
	}

	return nil
}
