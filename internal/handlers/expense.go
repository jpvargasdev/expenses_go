package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"guilliman/internal/models"
	"guilliman/internal/utils/timeutils"

	"github.com/gin-gonic/gin"
)

func AddExpenseHandler(c *gin.Context) {
	var mExpense models.Expense
	if err := c.ShouldBindJSON(&mExpense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	expense, err := models.AddExpense(mExpense)
	if err != nil {
		// You can log the error or return it, depending on your application's needs
		log.Printf("Error adding expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add expense"})
		return
	}
	c.JSON(http.StatusCreated, expense)
}

func GetExpensesHandler(c *gin.Context) {
	expenses, err := models.GetExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func GetExpensesForPeriodHandler(c *gin.Context) {
	dateParam := c.Query("date")
	var date time.Time
	if dateParam == "" {
		date = time.Now()
	} else {
		timestamp, err := strconv.ParseInt(dateParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use a Unix timestamp."})
			return
		}
		date = time.Unix(timestamp, 0)
	}

	start, end := timeutils.CalculatePeriodBoundaries(date)

	expenses, err := models.GetExpensesForPeriod(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func GetExpensesMonthlyHandler(c *gin.Context) {
	startDate, endDate := timeutils.GetSalaryMonthRange()
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	expenses, err := models.GetExpensesForPeriod(startTimestamp, endTimestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func RemoveExpenseHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	err = models.DeleteExpense(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})

}
