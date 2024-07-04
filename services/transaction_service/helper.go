package transaction_service

import (
	model "bank-grpc/model"
	"fmt"

	"gorm.io/gorm"
)

func checkAccExists(accountNum uint32, db *gorm.DB) (bool, error) {
	var acc []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&acc)
	if result.Error != nil {
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func authenticateAccount(accountNum uint32, accountPass string, db *gorm.DB) (bool, error) {
	var acc [] model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum, AccountPassword: accountPass}).Find(&acc)
	if result.Error != nil {
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func checkSufficientBalance(accountNum uint32, balance uint32, db *gorm.DB) (bool, error) {
	var accs []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&accs)
	if result.Error != nil {
		return false, fmt.Errorf("failed to fetch from db")
	}

	acc := accs[0]
	return balance <= acc.Balance, nil
}

func makeTransaction(fromAcc, toAcc, amt uint32, msg string, db *gorm.DB) error {

	// Debit in sender
	var accs []model.Account
	result := db.Where(&model.Account{AccountNumber: fromAcc}).Find(&accs)
	if result.Error != nil {
		return fmt.Errorf("failed to fetch from db")
	}

	acc := accs[0]
	acc.Balance -= amt
	db.Save(&acc)

	// Credit in receiver
	result = db.Where(&model.Account{AccountNumber: toAcc}).Find(&accs)
	if result.Error != nil {
		return fmt.Errorf("failed to fetch from db")
	}

	acc = accs[0]
	acc.Balance += amt
	db.Save(&acc)

	// Add new transaction
	transaction := model.Transaction{
		FromAccountNumber: fromAcc,
		ToAccountNumber: toAcc,
		Amount: amt,
		Message: msg,
	}
	db.Save(&transaction)

	return nil
}