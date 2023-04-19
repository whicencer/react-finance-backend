package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/handlers"
	"github.com/whicencer/react-finance-backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Auth Group
	auth := app.Group("/auth")

	auth.Get("/check", handlers.CheckToken)
	auth.Post("/signup", handlers.Register)
	auth.Post("/signin", handlers.Login)
	auth.Get("/me", middleware.AuthMiddleware, handlers.GetMe)

	// Me Group
	user := app.Group("/me", middleware.AuthMiddleware)

	// Cards
	user.Get("/cards", handlers.GetCards)
	user.Post("/cards", handlers.CreateCard)
	user.Delete("/cards", handlers.DeleteCard)
	user.Post("/cards/updateName", handlers.UpdateCardName)
	user.Post("/cards/updateTheme", handlers.UpdateCardTheme)

	// Transaction
	user.Post("/transactions", handlers.CreateTransaction)
	user.Get("/transactions", handlers.GetTransactions)
	user.Delete("/transactions", handlers.DeleteTransaction)
}
