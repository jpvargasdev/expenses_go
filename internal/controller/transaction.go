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

func (h *Controller) GetTransactionsController(c *gin.Context) {
	typeParam := c.Query("type")
	accountParam := c.Query("account")

	// check transaction type is valid or empty
	if typeParam != models.TransactionTypeExpense &&
		typeParam != models.TransactionTypeIncome &&
		typeParam != models.TransactionTypeSavings &&
		typeParam != models.TransactionTypeTransfer &&
		typeParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	transactions, err := models.GetTransactions(typeParam, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *Controller) AddTransactionController(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := models.AddTransaction(transaction)
	if err != nil {
		log.Printf("Error adding transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add transaction"})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

func (h *Controller) DeleteTransactionController(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = models.DeleteTransaction(id)
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
	dateParam := c.Query("date")
	typeParam := c.Query("type")
	accountParam := c.Query("account")

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

	expenses, err := models.GetTransactionsForPeriod(startTimestamp, endTimestamp, typeParam, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *Controller) GetTransactionsMonthlyController(c *gin.Context) {
	typeParam := c.Query("type")
	accountParam := c.Query("account")

	// check transaction type is valid or empty
	if typeParam != models.TransactionTypeExpense &&
		typeParam != models.TransactionTypeIncome &&
		typeParam != models.TransactionTypeSavings &&
		typeParam != models.TransactionTypeTransfer &&
		typeParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	startDate, endDate := timeutils.GetSalaryMonthRange()
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	expenses, err := models.GetTransactionsForPeriod(startTimestamp, endTimestamp, typeParam, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
