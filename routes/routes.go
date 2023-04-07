package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/handlers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.Root)
}
