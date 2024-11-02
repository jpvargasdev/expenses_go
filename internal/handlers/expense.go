package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"guilliman/internal/models"
	"guilliman/internal/utils/timeutils"

	"github.com/gin-gonic/gin"
)

func AddExpenseHandler(c *gin.Context) {
  var expense models.Expense
  if err := c.ShouldBindJSON(&expense); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  models.AddExpense(expense)
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

