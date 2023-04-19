package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/helpers"
	"github.com/whicencer/react-finance-backend/models"
)

const (
	Income  string = "income"
	Expense string = "expense"
)

// CreateTransaction
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
		return helpers.HandleBadRequest(c, "Invalid request body")
	}

	if body.Sum <= 0 {
		return helpers.HandleBadRequest(c, "Transaction sum can't beless than 1")
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

	// Increment or Decrement card balance
	var card models.Card
	if err := db.Where(models.Card{ID: body.BalanceId, UserID: uint(userId)}).First(&card).Error; err != nil {
		return helpers.HandleNotFound(c, "Card with specified ID not found")
	}

	if body.Status == Income {
		card.Balance += body.Sum
	} else if body.Status == Expense {
		if card.Balance < body.Sum {
			return helpers.HandleBadRequest(c, "Transaction expense sum can't be more than card balance")
		}
		card.Balance -= body.Sum
	} else {
		return helpers.HandleBadRequest(c, "Invalid transaction status")
	}

	if err := db.Save(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured on saving: "+err.Error())
	}

	// Creating Transaction
	if err := db.Create(&transaction).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message":     "Transaction has been successfully created",
		"ok":          true,
		"transaction": transaction,
	})
}

// GetTransactions
func GetTransactions(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)

	var transactions []models.Transaction
	if err := db.Where(&models.Transaction{UserID: uint(userId)}).Find(&transactions).Error; err != nil {
		return helpers.HandleNotFound(c, "Transactions were not found")
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"ok":           true,
	})
}

// DeleteTransaction
func DeleteTransaction(c *fiber.Ctx) error {
	db := database.DB
	claims := c.Locals("claims").(jwt.MapClaims)
	userId := claims["id"].(float64)
	transactionId := c.FormValue("transaction_id")

	if transactionId == "" {
		return helpers.HandleBadRequest(c, "Missing transaction id")
	}

	var transaction models.Transaction
	if err := db.Where(&models.Transaction{ID: transactionId, UserID: uint(userId)}).First(&transaction).Error; err != nil {
		return helpers.HandleNotFound(c, "Transaction was not found")
	}

	if err := db.Unscoped().Delete(&transaction).Error; err != nil {
		return helpers.HandleBadRequest(c, "Some error occured: "+err.Error())
	}

	var card models.Card
	if err := db.Where(models.Card{ID: transaction.BalanceId}).First(&card).Error; err != nil {
		return helpers.HandleNotFound(c, "Card with specified ID not found")
	}

	switch transaction.Status {
	case Income:
		card.Balance -= transaction.Sum
	case Expense:
		card.Balance += transaction.Sum
	default:
		return helpers.HandleBadRequest(c, "Invalid transaction status")
	}

	if err := db.Save(&card).Error; err != nil {
		return helpers.HandleInternalServerError(c, "Some error occured on saving: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Transaction was successfully deleted",
		"ok":      true,
	})
}
