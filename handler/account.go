package handler

import (
	"net/http"
	"task-golang-db/model"

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
