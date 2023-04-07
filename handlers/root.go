package handlers

import "github.com/gofiber/fiber/v2"

func Root(c *fiber.Ctx) error {
	return c.SendString("Hello world!")
}
