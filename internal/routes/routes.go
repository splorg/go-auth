package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/splorg/go-auth/internal/controller"
	"github.com/splorg/go-auth/internal/middleware"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", controller.HelloWorld)
	app.Post("/auth/register", controller.Register)
	app.Post("/auth/login", controller.Login)
	app.Get("/protected", middleware.Authenticate, controller.ProtectedRoute)
	app.Get("/user", controller.User)
}
