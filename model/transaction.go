package model

import (
	"bank-grpc/enums"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromAccountNumber uint32                `json:"from_account_number"`
	ToAccountNumber   uint32                `json:"to_account_number"`
	TransactionType   enums.TransactionType `json:"transaction_type"`
	Amount            uint32                `json:"amount"`
	Message           string                `json:"message"`
}
