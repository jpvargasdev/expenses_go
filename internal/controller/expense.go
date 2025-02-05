package controller

import (
	"net/http"

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

	expenses, err := models.GetTransactions(models.TransactionTypeExpense, accountParam, "", uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
