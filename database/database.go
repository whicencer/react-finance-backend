package database

import (
	"log"
	"os"

	"github.com/whicencer/react-finance-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to database, \n", err)
	}

	log.Println("Connected")

	db.AutoMigrate(&models.User{}, &models.Card{}, &models.Transaction{})

	DB = db
}
