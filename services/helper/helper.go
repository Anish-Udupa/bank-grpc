package helper

import (
	"bank-grpc/enums"
	model "bank-grpc/model"
	"fmt"
	"time"

	"gorm.io/datatypes"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func CheckAccExists(accountNum uint32, db *gorm.DB) (bool, error) {
	var acc []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&acc)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func AuthenticateAccount(accountNum uint32, accountPass string, db *gorm.DB) (bool, error) {
	var acc [] model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum, AccountPassword: accountPass}).Find(&acc)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return false, fmt.Errorf("failed to fetch from db")
	}
	return len(acc) > 0, nil
}

func CheckSufficientBalance(accountNum uint32, balance uint32, db *gorm.DB) (bool, error) {
	var accs []model.Account
	result := db.Where(&model.Account{AccountNumber: accountNum}).Find(&accs)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return false, fmt.Errorf("failed to fetch from db")
	}

	acc := accs[0]
	return balance <= acc.Balance, nil
}

func MakeTransaction(fromAcc, toAcc, amt uint32, msg string, db *gorm.DB) error {

	// Debit in sender
	var accs []model.Account
	result := db.Where(&model.Account{AccountNumber: fromAcc}).Find(&accs)
	if result.Error != nil {
		log.Errorf("failed to fetch from db: %v", result.Error)
		return fmt.Errorf("failed to fetch from db")
	}

	if len(accs) == 0 {
		log.Errorf("Account doesnt exist for accountNum: %v", fromAcc)
		return fmt.Errorf("account doesnt exist")
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
		return fmt.Errorf("account doesnt exist")
	}
	acc = accs[0]
	acc.Balance += amt
	db.Save(&acc)

	// Add new transaction
	transaction := model.Transaction{
		FromAccountNumber: fromAcc,
		ToAccountNumber: toAcc,
		TransactionType: enums.Transfer,
		Amount: amt,
		Message: msg,
	}
	db.Save(&transaction)

	return nil
}

func MakeDeposit(accountNum uint32, amount uint32, message string, db *gorm.DB) error {
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


	transaction := model.Transaction{
		FromAccountNumber: 0,
		ToAccountNumber: accountNum,
		TransactionType: enums.Deposit,
		Amount: amount,
		Message: message,
	}
	db.Save(&transaction)

	return nil
}

func MakeWithdraw(accountNum uint32, amount uint32, message string, db *gorm.DB) error {
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


	transaction := model.Transaction{
		FromAccountNumber: accountNum,
		ToAccountNumber: 0,
		TransactionType: enums.Withdraw,
		Amount: amount,
		Message: message,
	}

	db.Save(&transaction)
	return nil
}

func GetAccountStatements(accountNum uint32, fromDate, toDate time.Time, count uint32, db *gorm.DB) ([]model.Transaction, error) {
	var transactions []model.Transaction
	// Fetch from transaction
	res := db.Where("(from_account_number = ? OR to_account_number = ?) AND updated_at >= ? AND updated_at <= ?", accountNum, accountNum, datatypes.Date(fromDate), datatypes.Date(toDate)).Order("updated_at DESC, id").Find(&transactions)
	if res.Error != nil {
		return nil, res.Error
	}
	if count != 0 {
		transactions = transactions[:count]
	}

	return transactions, nil
}