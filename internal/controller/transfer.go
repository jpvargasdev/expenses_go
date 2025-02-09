package controller

import (
	"net/http"

	"guilliman/internal/models"
	"guilliman/internal/utils"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetTransfersController(c *gin.Context) {
	uid, err := utils.GetUserUID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Query("account")

	expenses, err := models.GetTransactions(models.TransactionTypeTransfer, id, "", uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func (h *Controller) TransferFundsController(c *gin.Context) {
	uid, err := utils.GetUserUID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transfer models.Transaction
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  if transfer.AccountID == transfer.RelatedAccountID {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Source and destination accounts cannot be the same"})
    return
  }

	transfer.UserID = uid

	// Validate required fields
	if !transfer.AccountID.Valid || !transfer.RelatedAccountID.Valid {
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
