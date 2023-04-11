package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/handlers"
	"github.com/whicencer/react-finance-backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Auth Group
	auth := app.Group("/auth")

	auth.Post("/signup", handlers.Register)
	auth.Post("/signin", handlers.Login)
	auth.Get("/me", middleware.AuthMiddleware, handlers.GetMe)

	// Me Group
	user := app.Group("/me", middleware.AuthMiddleware)

	user.Get("/cards", handlers.GetCards)
	user.Post("/cards", handlers.CreateCard)
	user.Post("/cards/updateName", handlers.UpdateCardName)
	user.Post("/cards/updateTheme", handlers.UpdateCardTheme)
}
