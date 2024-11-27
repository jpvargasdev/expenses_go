package controller

import (
	"net/http"
	"strconv"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetExpensesController(c *gin.Context) {
	accountParam := c.Query("account_id")
	accountId, _ := strconv.Atoi(accountParam)

	expenses, err := models.GetTransactions(models.TransactionTypeExpense, accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
