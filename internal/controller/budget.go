package controller

import (
	"guilliman/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)



func (h *Controller) GetBudgetSummaryController(c *gin.Context) {
  budgetSummary, err := models.GetBudgetSummary()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, budgetSummary)
}
