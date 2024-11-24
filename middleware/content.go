package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AllowedContentType(allowedContentTypes []string) fiber.Handler {
	if len(allowedContentTypes) == 0 {
		panic("allowedContentTypes is empty")
	}

	return func(ctx *fiber.Ctx) error {
		var reqContentType string

		if headers, ok := ctx.GetReqHeaders()[fiber.HeaderContentType]; !ok || len(headers) == 0 {
			return fiber.ErrBadRequest
		} else {
			reqContentType = headers[0]
		}

		for _, allowedContentType := range allowedContentTypes {
			if reqContentType == allowedContentType {
				return ctx.Next()
			}
		}

		return fiber.ErrUnsupportedMediaType
	}
}

func AllowedContentTypeWithJSON() fiber.Handler {
	return AllowedContentType([]string{fiber.MIMEApplicationJSON})
}
