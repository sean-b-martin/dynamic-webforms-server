package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/auth"
	"strings"
)

var UserIDLocal = "userID"

type JWTAuth struct {
	jwtService *auth.JWTService
}

func NewJWTAuth(jwtService *auth.JWTService) *JWTAuth {
	return &JWTAuth{jwtService}
}

func (j *JWTAuth) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header, ok := c.GetReqHeaders()[fiber.HeaderAuthorization]
		if !ok || len(header) == 0 {
			return fiber.ErrUnauthorized
		}

		var token string
		if token, ok = strings.CutPrefix(header[0], "Bearer "); !ok {
			return fiber.ErrUnauthorized
		}

		claims, err := j.jwtService.ValidateToken(token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals(UserIDLocal, userID)

		return c.Next()
	}
}
