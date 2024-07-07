package transaction_service

import (
	model "bank-grpc/model"
	"fmt"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func checkAccExists(accountNum uint32, db *gorm.DB) (bool, error) {
	var acc []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&acc)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func authenticateAccount(accountNum uint32, accountPass string, db *gorm.DB) (bool, error) {
	var acc [] model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum, AccountPassword: accountPass}).Find(&acc)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func checkSufficientBalance(accountNum uint32, balance uint32, db *gorm.DB) (bool, error) {
	var accs []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&accs)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
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
		log.Errorf("failed to fetch from db: %v", result.Error)
		return fmt.Errorf("failed to fetch from db")
	}

	if len(accs) == 0 {
		log.Errorf("Account doesnt exist for accountNum: %v", fromAcc)
		return fmt.Errorf("Account doesnt exist")
	}
	acc := accs[0]
	acc.Balance -= amt
	db.Save(&acc)

	// Credit in receiver
	result = db.Where(&model.Account{AccountNumber: toAcc}).Find(&accs)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return fmt.Errorf("failed to fetch from db")
	}

	if len(accs) == 0 {
		log.Errorf("Account doesnt exist for accountNum: %v", toAcc)
		return fmt.Errorf("Account doesnt exist")
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

func makeDeposit(accountNum uint32, amount uint32, message string, db *gorm.DB) error {
	// Credit to account
	var accounts []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&accounts)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return fmt.Errorf("failed to fetch from db")
	}

	if len(accounts) == 0 {
		log.Errorf("Account doesnt exist for accountNum: %v", accountNum)
		return fmt.Errorf("account doesnt exist")
	}

	acc := accounts[0]
	acc.Balance += amount
	db.Save(&acc)

	deposit := model.DepositWithdraw{
		AccountNumber: accountNum,
		Amount: amount,
		Message: message,
		Type: "DEPOSIT",
	}
	db.Save(&deposit)
	return nil
}

func makeWithdraw(accountNum uint32, amount uint32, message string, db *gorm.DB) error {
	// Debit from account
	var accounts []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&accounts)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return fmt.Errorf("failed to fetch from db")
	}

	if len(accounts) == 0 {
		log.Errorf("Account doesnt exist for accountNum: %v", accountNum)
		return fmt.Errorf("account doesnt exist")
	}

	acc := accounts[0]
	acc.Balance -= amount
	db.Save(&acc)

	withdraw := model.DepositWithdraw{
		AccountNumber: accountNum,
		Amount: amount,
		Message: message,
		Type: "WITHDRAW",
	}
	db.Save(&withdraw)
	return nil
}