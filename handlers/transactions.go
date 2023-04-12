package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/models"
)

func CreateTransaction(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)

	var body struct {
		Status    string
		Sum       int
		Note      string
		BalanceId string
		Category  string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"ok":      false,
		})
	}

	if body.Sum <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Transaction sum can't beless than 1",
			"ok":      false,
		})
	}

	transaction := models.Transaction{
		ID:        uuid.New().String(),
		BalanceId: body.BalanceId,
		Category:  body.Category,
		Note:      body.Note,
		Status:    body.Status,
		Sum:       body.Sum,
		UserID:    uint(userId),
	}

	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Transaction has been successfully created",
		"ok":          true,
		"transaction": transaction,
	})
}
