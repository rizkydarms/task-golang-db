package model

type Transaction struct {
	TransactionID         int8   `json:"transaction_id" gorm:"primaryKey;autoIncrement;<-:false"`
	TransactionCategoryID int8   `json:"transaction_category_id"`
	AccountID             int8   `json:"account_id"`
	FromAccountId         int8   `json:"from_account_id"`
	ToAccountId           int8   `json:"to_account_id"`
	Amount                int8   `json:"amount"`
	TransactionDate       string `json:"transaction_date"`
}
