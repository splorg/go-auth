package main

import (
	"github.com/splorg/go-auth/internal/database"
	"github.com/splorg/go-auth/internal/routes"
	"github.com/splorg/go-auth/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	validator.Setup()
	app := fiber.New()
	app.Use(cors.New(cors.Config{}))

	routes.RegisterRoutes(app)

	app.Listen(":3000")
}
