package controller

import (
	"net/http"
	"strconv"

	"guilliman/internal/models"
	"guilliman/internal/utils"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetExpensesController(c *gin.Context) {
	uid, err := utils.GetUserUID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accountParam := c.Query("account_id")
	accountId, _ := strconv.Atoi(accountParam)

	expenses, err := models.GetTransactions(models.TransactionTypeExpense, accountId, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
