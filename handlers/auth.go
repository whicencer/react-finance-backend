package handlers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/whicencer/react-finance-backend/database"
	"github.com/whicencer/react-finance-backend/helpers"
	"github.com/whicencer/react-finance-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// Register
func Register(c *fiber.Ctx) error {
	db := database.DB

	var body struct {
		Username string
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		helpers.HandleBadRequest(c, "Invalid request body")
	}

	if len(body.Password) < 8 {
		helpers.HandleBadRequest(c, "Password length should be 8 symbols or more")
	}

	if len(body.Username) < 2 {
		helpers.HandleBadRequest(c, "Username length should be 2 symbols or more")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		helpers.HandleInternalServerError(c, "Failed to hash password")
	}

	user := models.User{
		Username: body.Username,
		Password: string(hashedPassword),
	}

	var existingUser models.User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return helpers.HandleBadRequest(c, "This username is already taken")
	}

	if err := db.Create(&user).Error; err != nil {
		helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
	}

	card := models.Card{
		ID:       uuid.New().String(),
		Balance:  0,
		CardName: "General",
		ThemeId:  1,
		UserID:   user.ID,
	}

	if err := db.Create(&card).Error; err != nil {
		helpers.HandleInternalServerError(c, "Some error occured: "+err.Error())
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
		helpers.HandleBadRequest(c, "Invalid request body")
	}

	var dbUser models.User

	if err := db.Where("username = ?", body.Username).First(&dbUser).Error; err != nil {
		helpers.HandleUnauthorized(c, "Invalid login or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(body.Password)); err != nil {
		helpers.HandleUnauthorized(c, "Invalid login or password")
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
		helpers.HandleNotFound(c, "Cannot find user")
	}

	return c.JSON(fiber.Map{
		"user": user,
		"ok":   true,
	})
}
