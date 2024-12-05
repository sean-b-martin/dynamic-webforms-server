package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/validation"
)

// TODO change function to only return one value
func parseAndValidateRequestData(ctx *fiber.Ctx, paramsOut interface{}, bodyOut interface{}) (error, bool) {
	if paramsOut != nil {
		if err := ctx.ParamsParser(paramsOut); err != nil {
			return err, false
		}

		if validationErrors := validation.Validate(paramsOut); len(validationErrors) > 0 {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors), false
		}
	}

	if bodyOut != nil {
		if err := ctx.BodyParser(bodyOut); err != nil {
			return err, false
		}

		if validationErrors := validation.Validate(bodyOut); len(validationErrors) > 0 {
			ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors), false
		}
	}

	return nil, true
}
