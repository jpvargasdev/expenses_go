package controller

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

func (h *Controller) GetTransactionsByAccountController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	accountId := c.Param("id")
	transactions, err := models.GetTransactionsByAccount(accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *Controller) GetTransactionsByMainCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	mainCategory := c.Param("main_category")
	startDay := c.Query("start_day")
	endDay := c.Query("end_day")

	transactions, err := models.GetTransactionsByMainCategory(mainCategory, startDay, endDay, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Controller) GetTransactionsController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	typeParam := c.Query("type")
	accountParam := c.Query("account")
	accountId, _ := strconv.Atoi(accountParam)

	// check transaction type is valid or empty
	if typeParam != models.TransactionTypeExpense &&
		typeParam != models.TransactionTypeIncome &&
		typeParam != models.TransactionTypeSavings &&
		typeParam != models.TransactionTypeTransfer &&
		typeParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	transactions, err := models.GetTransactions(typeParam, accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *Controller) AddTransactionController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := models.AddTransaction(transaction, userID.(int))
	if err != nil {
		log.Printf("Error adding transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transaction"})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

func (h *Controller) DeleteTransactionController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = models.DeleteTransaction(id, userID.(int))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func (h *Controller) GetTransactionsForPeriodController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	dateParam := c.Query("date")
	typeParam := c.Query("type")
	accountParam := c.Query("account")
	accountId, _ := strconv.Atoi(accountParam)

	// check transaction type is valid or empty
	if typeParam != models.TransactionTypeExpense &&
		typeParam != models.TransactionTypeIncome &&
		typeParam != models.TransactionTypeSavings &&
		typeParam != models.TransactionTypeTransfer &&
		typeParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

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

	startTimestamp, endTimestamp := timeutils.CalculatePeriodBoundaries(date)

	expenses, err := models.GetTransactionsForPeriod(startTimestamp, endTimestamp, typeParam, accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *Controller) GetTransactionsMonthlyController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	typeParam := c.Query("type")
	accountParam := c.Query("account")
	accountId, _ := strconv.Atoi(accountParam)
	startDay := c.Query("start_day")
	endDay := c.Query("end_day")

	// check transaction type is valid or empty
	if typeParam != models.TransactionTypeExpense &&
		typeParam != models.TransactionTypeIncome &&
		typeParam != models.TransactionTypeSavings &&
		typeParam != models.TransactionTypeTransfer &&
		typeParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	startDate, endDate := timeutils.GetSalaryMonthRange(startDay, endDay)
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	expenses, err := models.GetTransactionsForPeriod(startTimestamp, endTimestamp, typeParam, accountId, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
