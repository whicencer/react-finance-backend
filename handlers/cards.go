package handlers

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/helpers"
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
		return helpers.HandleBadRequest(c, "Invalid request body")
	}

	if body.Balance <= 0 {
		return helpers.HandleBadRequest(c, "Balance can't be less than 1")
	}

	card := models.Card{
		ID:       uuid.New().String(),
		Balance:  body.Balance,
		CardName: body.CardName,
		ThemeId:  1,
		UserID:   uint(userId),
	}

	if err := db.Create(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
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

	if cardId == "" || newName == "" {
		return helpers.HandleBadRequest(c, "Invalid request body")
	}

	var card models.Card

	if err := db.Where(&models.Card{ID: cardId, UserID: uint(userId)}).First(&card).Error; err != nil {
		return helpers.HandleNotFound(c, "Card was not found")
	}

	card.CardName = newName

	if err := db.Save(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Card name was successfully changed",
		"ok":      true,
	})
}

func UpdateCardTheme(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)

	userId := claims["id"].(float64)
	cardId := c.FormValue("card_id")
	newThemeIdStr := c.FormValue("theme_id")

	if cardId == "" || newThemeIdStr == "" {
		return helpers.HandleBadRequest(c, "Invalid request body")
	}

	newThemeId, err := strconv.Atoi(newThemeIdStr)

	if err != nil {
		return helpers.HandleBadRequest(c, "theme_id should be a number")
	}

	var card models.Card

	if err := db.Where(&models.Card{ID: cardId, UserID: uint(userId)}).First(&card).Error; err != nil {
		return helpers.HandleNotFound(c, "Card was not found")
	}

	card.ThemeId = newThemeId

	if err := db.Save(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Card theme was successfully changed",
		"ok":      true,
	})

}

func DeleteCard(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)
	cardId := c.FormValue("card_id")

	if cardId == "" {
		return helpers.HandleBadRequest(c, "Invalid request body")
	}

	var card models.Card

	if err := db.Where(&models.Card{ID: cardId, UserID: uint(userId)}).First(&card).Error; err != nil {
		return helpers.HandleNotFound(c, "Card was not found")
	}

	if err := db.Unscoped().Delete(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
	}

	// Delete all transactions
	if err := db.Where(&models.Transaction{BalanceId: cardId}).Delete(&models.Transaction{}); err != nil {
		return helpers.HandleBadRequest(c, "Some error occured on deleting transactions")
	}

	return c.JSON(fiber.Map{
		"message": "Card was successfully deleted",
		"ok":      true,
	})
}
