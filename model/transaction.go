package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	FromAccountNumber uint32
	ToAccountNumber   uint32
	Amount            uint32
	Message           string
}
