package helpers

import "github.com/gofiber/fiber/v2"

func HandleBadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message,
		"ok":      false,
	})
}

func HandleNotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": message,
		"ok":      false,
	})
}

func HandleInternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": message,
		"ok":      false,
	})
}

func HandleUnauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": message,
		"ok":      false,
	})
}
