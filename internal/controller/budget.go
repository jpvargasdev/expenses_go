package controller

import (
	"guilliman/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetBudgetSummaryController(c *gin.Context) {
	startDay := c.Query("start_day")
	endDay := c.Query("end_day")

	budgetSummary, err := models.GetBudgetSummary(startDay, endDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgetSummary)
}
