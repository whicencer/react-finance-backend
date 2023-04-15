package handlers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "This username is already taken",
				"ok":      false,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Some error occured: " + err.Error(),
				"ok":      false,
			})
		}
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
	db := database.DB

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"ok":      false,
		})
	}

	var dbUser models.User

	if err := db.Where("username = ?", body.Username).First(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid login or password",
			"ok":      false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid login or password",
			"ok":      false,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = dbUser.Username
	claims["id"] = dbUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "You have successfully logged in",
		"token":   t,
		"ok":      true,
	})
}

// Get me
func GetMe(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	claims := c.Locals("claims").(jwt.MapClaims)
	id := claims["id"].(float64)

	if err := db.Where("ID = ?", id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Cannot find user",
			"ok":      false,
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
		"ok":   true,
	})
}
