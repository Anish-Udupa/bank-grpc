package enums

type TransactionType string

const (
	Transfer TransactionType = "TRANSFER"
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)
