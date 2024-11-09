package handler

import (
	"net/http"
	"task-golang-db/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewTransactionInterface interface {
	NewTransaction(*gin.Context)
	TransactionList(*gin.Context)
}

type newTransactionImplement struct {
	db *gorm.DB
}

func NewTrans(db *gorm.DB) NewTransactionInterface {
	return &newTransactionImplement{
		db: db,
	}
}

// NewTransaction creates a new transaction and updates the account balance
func (a *newTransactionImplement) NewTransaction(c *gin.Context) {
	var data struct {
		AccountID             int8  `json:"account_id"`
		TransactionCategoryID *int8 `json:"transaction_category_id"`
		FromAccountId         *int8 `json:"from_account_id"`
		ToAccountId           *int8 `json:"to_account_id"`
		Amount                int8  `json:"amount"`
	}

	// Bind JSON to data struct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new transaction record
	transaction := model.Transaction{
		AccountID: data.AccountID,
		Amount:    data.Amount,
	}
	if data.TransactionCategoryID != nil {
		transaction.TransactionCategoryID = *data.TransactionCategoryID
	}
	if data.FromAccountId != nil {
		transaction.FromAccountId = *data.FromAccountId
	}
	if data.ToAccountId != nil {
		transaction.ToAccountId = *data.ToAccountId
	}
	// Save transaction to the database
	if err := a.db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the account and update balance
	var account model.Account
	if err := a.db.First(&account, data.AccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	account.Balance += float64(data.Amount)
	a.db.Save(&account)

	// Return the created transaction
	c.JSON(http.StatusOK, transaction)
}

// TransactionList retrieves transactions by account_id, ordered by transaction date
func (a *newTransactionImplement) TransactionList(c *gin.Context) {
	accountID := c.Query("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	var transaction []model.Transaction
	if err := a.db.Where("account_id = ?", accountID).Order("transaction_date desc").Find(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
