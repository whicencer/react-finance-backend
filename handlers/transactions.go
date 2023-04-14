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

	const (
		Income  string = "income"
		Expense string = "expense"
	)

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

	var card models.Card
	if err := db.Where(models.Card{ID: body.BalanceId}).First(&card).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Card with specified ID not found",
			"ok":      false,
		})
	}

	if body.Status == Income {
		card.Balance += body.Sum
	} else if body.Status == Expense {
		if card.Balance < body.Sum {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Transaction expense sum can't be more than card balance",
				"ok":      false,
			})
		}
		card.Balance -= body.Sum
	}
	if err := db.Save(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured on saving: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Transaction has been successfully created",
		"ok":          true,
		"transaction": transaction,
	})
}

func GetTransactions(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)

	var transactions []models.Transaction
	if err := db.Where(&models.Transaction{UserID: uint(userId)}).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Transactions were not found",
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"ok":           true,
	})
}

func DeleteTransaction(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)
	transactionId := c.FormValue("transaction_id")

	if transactionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing transaction id",
			"ok":      false,
		})
	}

	var transaction models.Transaction
	if err := db.Where(&models.Transaction{ID: transactionId, UserID: uint(userId)}).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Transaction was not found",
			"ok":      false,
		})
	}

	if err := db.Unscoped().Delete(&transaction).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction was successfully deleted",
		"ok":      true,
	})
}
