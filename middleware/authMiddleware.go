package middleware

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if len(authHeader) <= len("Bearer ") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Auth token is missing",
			"ok":      false,
		})
	}

	authToken := authHeader[len("Bearer "):]

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"ok":      false,
		})
	}

	c.Locals("claims", token.Claims)

	return c.Next()
}
