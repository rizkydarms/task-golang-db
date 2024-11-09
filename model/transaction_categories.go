package model

type TransCat struct {
	AccountID             int64  `json:"account_id" gorm:"primaryKey;autoIncrement;<-:false"`
	TransactionCategoryID int64  `json:"transaction_category_id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name                  string `json:"name"`
}

// Untuk memastikan ORM menggunakan nama tabel yang benar
func (TransCat) TableName() string {
	return "transaction_categories"
}
