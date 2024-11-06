package model

type Account struct {
	AccountID int64   `json:"account_id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
}

// Jika ingin menggunakan nama tabel khusus
func (Account) TableName() string {
	return "accounts"
}
