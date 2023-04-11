package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ID        string `json:"id"`
	BalanceId string `json:"balanceId"`
	Category  string `json:"category"`
	Note      string `json:"note"`
	Status    string `json:"status"`
	Sum       int    `json:"sum"`
	UserID    uint   `json:"user_id"`
}
