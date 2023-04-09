package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/models"
)

func GetCards(c *fiber.Ctx) error {
	db := database.DB

	var cards []models.Card

	db.Where("user_id = ?", 8).Find(&cards)

	return c.JSON(fiber.Map{
		"cards": cards,
		"ok":    true,
	})
}
