package controller

import (
	"net/http"
	"strconv"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetIncomesController(c *gin.Context) {
	accountParam := c.Query("account")
	accountId, _ := strconv.Atoi(accountParam)

	incomes, err := models.GetTransactions(models.TransactionTypeIncome, accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incomes)
}
