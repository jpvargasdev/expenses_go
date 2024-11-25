package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func AddIncomeHandler(c *gin.Context) {
  var income models.Income
  if err := c.ShouldBindJSON(&income); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  models.AddIncome(income)
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

