package model

type Account struct {
	AccountID int64  `json:"account_id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name      string `json:"name"`
	Balance   string `json:"balance"`
}

// func (Account) TableName() string {
// 	return "accounts"
// }
