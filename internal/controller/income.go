package controller

import (
	"net/http"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetIncomesController(c *gin.Context) {
	accountParam := c.Query("account")

	incomes, err := models.GetTransactions(models.TransactionTypeIncome, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incomes)
}
