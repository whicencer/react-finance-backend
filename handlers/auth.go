package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// Register
func Register(c *fiber.Ctx) error {
	db := database.DB

	var body struct {
		Username string
		Balance  int
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"ok":      false,
		})
	}

	if len(body.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password length should be 8 symbols or more",
			"ok":      false,
		})
	}

	if len(body.Username) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username length should be 2 symbols or more",
			"ok":      false,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"ok":      false,
		})
	}

	user := models.User{
		Username: body.Username,
		Balance:  body.Balance,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	card := models.Card{
		ID:       uuid.New().String(),
		Balance:  0,
		CardName: "General",
		ThemeId:  1,
		UserID:   user.ID,
	}

	if err := db.Create(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some error occured: " + err.Error(),
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"user":    user,
		"card":    card,
		"message": "User created",
		"ok":      true,
	})
}

// Login
func Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

// Get me
func GetMe(c *fiber.Ctx) error {
	return c.SendString("Get me")
}
