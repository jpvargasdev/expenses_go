package controller

import (
	"net/http"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetTransfersController(c *gin.Context) {
	accountParam := c.Query("account")

	expenses, err := models.GetTransactions(models.TransactionTypeTransfer, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func (h *Controller) TransferFundsController(c *gin.Context) {
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

	transaction, err := models.AddTransfer(transfer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}
