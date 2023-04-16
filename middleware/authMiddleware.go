package middleware

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/helpers"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if len(authHeader) <= len("Bearer ") {
		return helpers.HandleBadRequest(c, "Auth token is missing")
	}

	authToken := authHeader[len("Bearer "):]

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return helpers.HandleUnauthorized(c, "Invalid token")
	}

	c.Locals("claims", token.Claims)

	return c.Next()
}
