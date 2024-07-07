package model

import "gorm.io/gorm"

type DepositWithdraw struct {
	gorm.Model
	AccountNumber uint32
	Amount        uint32
	Type          string
	Message       string
}