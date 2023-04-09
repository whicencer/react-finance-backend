package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/handlers"
)

func SetupRoutes(app *fiber.App) {
	// Auth Group
	auth := app.Group("/auth")

	auth.Post("/signup", handlers.Register)
	auth.Post("/signin", handlers.Login)
	auth.Get("/me", handlers.GetMe)

	// Me Group
	user := app.Group("/me")

	user.Get("/cards", handlers.GetCards)
}
