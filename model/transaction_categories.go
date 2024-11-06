package model

type TransCat struct {
	TransactionCategoryID int64  `json:"transaction_category_id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name                  string `json:"name"`
}

// Untuk memastikan ORM menggunakan nama tabel yang benar
func (TransCat) TableName() string {
	return "transaction_categories"
}
