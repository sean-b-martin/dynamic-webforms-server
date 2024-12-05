package controller

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

func serviceErrToResponse(ctx *fiber.Ctx, err error) error {
	if err != nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	if errors.Is(err, service.ErrNoPermission) {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	return ctx.SendStatus(fiber.StatusInternalServerError)
}
