package controller

import (
	"net/http"
	"strconv"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetExpensesController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	accountParam := c.Query("account_id")
	accountId, _ := strconv.Atoi(accountParam)

	expenses, err := models.GetTransactions(models.TransactionTypeExpense, accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
