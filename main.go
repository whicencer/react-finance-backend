package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/routes"
)

func main() {
	app := fiber.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	routes.SetupRoutes(app)

	log.Fatal(app.Listen("127.0.0.1:" + port))
}
