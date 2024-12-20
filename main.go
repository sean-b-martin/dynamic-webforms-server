package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sean-b-martin/dynamic-webforms-server/auth"
	"github.com/sean-b-martin/dynamic-webforms-server/controller"
	"github.com/sean-b-martin/dynamic-webforms-server/database"
	"github.com/sean-b-martin/dynamic-webforms-server/middleware"
	"github.com/sean-b-martin/dynamic-webforms-server/service"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// setup database
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config.json: %w", err))
	}

	var dbConfig database.ConnectionConfig
	if err = json.Unmarshal(data, &dbConfig); err != nil {
		log.Fatal(fmt.Errorf("error parsing config.json: %w", err))
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(dbConfig); err != nil {
		log.Fatal(fmt.Errorf("error parsing config.json: %w", err))
	}

	db, err := database.CreateDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	database.CreateTables(db)

	// setup webserver
	app := fiber.New(fiber.Config{
		JSONDecoder: func(data []byte, v interface{}) error {
			decoder := json.NewDecoder(bytes.NewReader(data))
			decoder.DisallowUnknownFields()
			return decoder.Decode(v)
		},
	})
	app.Use(recover.New())

	jwtService, err := auth.NewJWTService()
	if err != nil {
		log.Fatal(fmt.Errorf("error creating JWT service: %w", err))
	}

	passwordService, err := auth.NewPasswordService(bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating password service: %w", err))
	}

	authMiddleware := middleware.NewJWTAuth(jwtService)
	controller.NewUserController(app.Group("/users"), authMiddleware,
		service.NewUserService(db, passwordService, jwtService))
	controller.NewFormController(app.Group("/forms"), authMiddleware, service.NewFormService(db))
	controller.NewSchemaController(app.Group("/forms/:formID/"), authMiddleware, service.NewSchemaService(db))

	// shutdown server gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := app.ShutdownWithTimeout(1 * time.Minute); err != nil {
			log.Fatal(fmt.Errorf("error shutting down server: %w", err))
		}
	}()

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
