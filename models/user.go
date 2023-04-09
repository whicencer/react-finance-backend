package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Balance  int    `json:"balance"`
	Cards    []Card `json:"cards" gorm:"foreignKey:UserID"`
	Password string `json:"-"`
}

type Card struct {
	gorm.Model
	ID       string
	Balance  int    `json:"balance"`
	CardName string `json:"cardName"`
	ThemeId  int    `json:"themeId"`
	UserID   uint   `json:"user_id"`
}
