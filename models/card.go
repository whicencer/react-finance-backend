package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	ID       string `json:"card_id"`
	Balance  int    `json:"balance"`
	CardName string `json:"cardName"`
	ThemeId  int    `json:"themeId"`
	UserID   uint   `json:"user_id"`
}
