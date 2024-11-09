package handler

import (
	"net/http"
	"task-golang-db/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountInterface interface {
	Create(*gin.Context)
	Read(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	List(*gin.Context)
	My(*gin.Context)
	TopUp(*gin.Context)
	Balance(*gin.Context)
	Transfer(*gin.Context)
	Mutation(*gin.Context)
}

type accountImplement struct {
	db *gorm.DB
}

// Constructor untuk accountImplement
func NewAccount(db *gorm.DB) AccountInterface {
	return &accountImplement{
		db: db,
	}
}

// Implementasi metode Create (contoh implementasi)
func (a *accountImplement) Create(c *gin.Context) {
	var request model.Account
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.db.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully",
		"account": request,
	})
}

// Implementasi metode Read
func (a *accountImplement) Read(c *gin.Context) {
	accountID := c.Param("id")
	var account model.Account
	if err := a.db.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": account,
	})
}

// Implementasi metode Update
func (a *accountImplement) Update(c *gin.Context) {
	accountID := c.Param("id")
	var account model.Account
	if err := a.db.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var request model.Account
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.Name = request.Name
	account.Balance = request.Balance
	if err := a.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

// Implementasi metode Delete
func (a *accountImplement) Delete(c *gin.Context) {
	accountID := c.Param("id")
	if err := a.db.Delete(&model.Account{}, accountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// Implementasi metode List
func (a *accountImplement) List(c *gin.Context) {
	var accounts []model.Account
	if err := a.db.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// Implementasi metode My (menampilkan akun milik pengguna yang sedang login)
func (a *accountImplement) My(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID not provided"})
		return
	}

	var account model.Account
	if err := a.db.First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// Implementasi metode TopUp
func (a *accountImplement) TopUp(c *gin.Context) {
	var request struct {
		AccountID int64   `json:"account_id" binding:"required"`
		Amount    float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account model.Account
	if err := a.db.First(&account, request.AccountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	account.Balance += request.Amount
	if err := a.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Top-up successful",
		"balance": account.Balance,
	})
}

// Implementasi metode Balance
func (a *accountImplement) Balance(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID not provided"})
		return
	}

	var account model.Account
	if err := a.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

// Implementasi metode Transfer
func (a *accountImplement) Transfer(c *gin.Context) {
	AccountID := c.GetInt64("account_id")
	payload := struct {
		ToAccountID int64 `json:"to_account_id"`
		Amount      int64 `json:"amount"`
	}{}

	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the current and target accounts
	var senderAccount model.Account
	var receiverAccount model.Account

	if err := a.db.First(&senderAccount, AccountID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Sender account not found"})
		return
	}

	if err := a.db.First(&receiverAccount, payload.ToAccountID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Target account not found"})
		return
	}

	// Check balance and update if sufficient
	if senderAccount.Balance < float64(payload.Amount) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	senderAccount.Balance -= float64(payload.Amount)
	receiverAccount.Balance += float64(payload.Amount)

	if err := a.db.Save(&senderAccount).Error; err != nil || a.db.Save(&receiverAccount).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer balance"})
		return
	}

	// Convert AccountID and ToAccountID to int8
	fromAccountID := int8(AccountID)
	toAccountID := int8(payload.ToAccountID)
	amount := int8(payload.Amount)

	// Create transaction record
	transaction := model.Transaction{
		AccountID:       fromAccountID,
		FromAccountId:   fromAccountID,
		ToAccountId:     toAccountID,
		Amount:          amount,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"), // format sebagai string
	}
	if err := a.db.Create(&transaction).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

// Imp;ementasi metode Mutations
// Mutation returns a list of transactions for the current user, sorted by latest (requires auth)
func (a *accountImplement) Mutation(c *gin.Context) {
	accountID := c.GetInt64("account_id")

	var transactions []model.Transaction
	query := a.db.Where("account_id = ?", accountID).Order("transaction_date DESC")

	if err := query.Find(&transactions).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
