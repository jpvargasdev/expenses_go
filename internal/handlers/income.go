package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"guilliman/internal/models"
	"guilliman/internal/utils/timeutils"

	"github.com/gin-gonic/gin"
)

func AddIncomeHandler(c *gin.Context) {
	var income models.Income
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	income, err := models.AddIncome(income)
	if err != nil {
		// You can log the error or return it, depending on your application's needs
		log.Printf("Error adding expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add expense"})
		return
	}
	c.JSON(http.StatusCreated, income)
}

func GetIncomesHandler(c *gin.Context) {
	incomes, err := models.GetIncomes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incomes)
}

func GetIncomesMonthlyHandler(c *gin.Context) {
	startDate, endDate := timeutils.GetSalaryMonthRange()
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	incomes, err := models.GetIncomesForPeriod(startTimestamp, endTimestamp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incomes)
}

func RemoveIncomeHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid income ID"})
		return
	}

	err = models.DeleteIncome(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income deleted successfully"})

}
