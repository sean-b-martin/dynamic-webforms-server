package main

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sean-b-martin/dynamic-webforms-server/controller"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONDecoder: func(data []byte, v interface{}) error {
			decoder := json.NewDecoder(bytes.NewReader(data))
			decoder.DisallowUnknownFields()
			return decoder.Decode(v)
		},
	})

	app.Use(recover.New())
	controller.NewUserController(app.Group("/users"))
	log.Fatal(app.Listen(":3000"))
}
