package controller

import (
	"net/http"
	"strconv"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetTransfersController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	accountParam := c.Query("account")
	accountId, _ := strconv.Atoi(accountParam)

	expenses, err := models.GetTransactions(models.TransactionTypeTransfer, accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func (h *Controller) TransferFundsController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	var transfer models.Transaction
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if transfer.AccountID == 0 || transfer.RelatedAccountID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both source and destination accounts are required"})
		return
	}
	if transfer.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transfer amount must be greater than zero"})
		return
	}

	transaction, err := models.AddTransfer(transfer, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}
