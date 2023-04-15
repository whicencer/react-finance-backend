package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/routes"
)

func main() {
	app := fiber.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// CORS
	app.Use(cors.New())

	// routes
	routes.SetupRoutes(app)

	// database
	database.Connect()

	// listening
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
