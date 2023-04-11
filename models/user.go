package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string        `json:"username" gorm:"unique"`
	Balance      int           `json:"balance"`
	Cards        []Card        `json:"cards" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID"`
	Password     string        `json:"-"`
}
