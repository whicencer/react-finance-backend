package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/models"
)

func GetMyCards(c *fiber.Ctx) error {
	db := database.DB
	var cards []models.Card

	claims := c.Locals("claims").(jwt.MapClaims)

	id := claims["id"].(float64)

	db.Where("user_id = ?", id).Find(&cards)

	return c.JSON(fiber.Map{
		"cards": cards,
		"ok":    true,
	})
}
