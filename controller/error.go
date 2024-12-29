package controller

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
)

func handleServiceErr(ctx *fiber.Ctx, err error) error {
	if err != nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	if errors.Is(err, service.ErrNoPermission) {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	log.Error(err.Error())
	return ctx.SendStatus(fiber.StatusInternalServerError)
}
