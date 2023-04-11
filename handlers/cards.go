package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/models"
)

func GetCards(c *fiber.Ctx) error {
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

func CreateCard(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)

	var body struct {
		CardName string
		Balance  int
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"ok":      false,
		})
	}

	if body.Balance <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Balance can't be less than 1",
			"ok":      false,
		})
	}

	card := models.Card{
		ID:       uuid.New().String(),
		Balance:  body.Balance,
		CardName: body.CardName,
		ThemeId:  1,
		UserID:   uint(userId),
	}

	if err := db.Create(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"card":    card,
		"message": "Card has been successfully created",
		"ok":      true,
	})
}

func UpdateCardName(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)

	cardId := c.FormValue("card_id")
	newName := c.FormValue("card_name")

	var card models.Card

	if err := db.Where(&models.Card{ID: cardId, UserID: uint(userId)}).First(&card).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Card was not found",
			"ok":      false,
		})
	}

	card.CardName = newName

	if err := db.Save(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Card name was successfully changed",
		"ok":      true,
	})
}
