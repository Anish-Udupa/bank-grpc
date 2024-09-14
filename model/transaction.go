package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	FromAccountNumber uint32 `json:"from_account_number"`
	ToAccountNumber   uint32
	Amount            uint32
	Message           string
}
