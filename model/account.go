package model

type Account struct {
	AccountNumber     uint32 `json:"account_number,omitempty" gorm:"primaryKey"`
	AccountHolderName string `json:"account_holder_name,omitempty"`
	AccountPassword   string `json:"account_password,omitempty"`
	Balance           uint32 `json:"balance,omitempty"`
}
