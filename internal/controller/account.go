package controller

import (
	"guilliman/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAccounts godoc
// @Summary      Get accounts
// @Description  get accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]

func (h *Controller) GetAccountsController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	accountId := c.Param("id")

	accounts, err := models.GetAccounts(accountId, userID.(int)) // Fetch accounts from storage
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *Controller) AddAccountController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var newAccount models.Account
	if err := c.ShouldBindJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := models.AddAccount(newAccount, userID.(int)) // Add account to storage
	if err != nil {
		// You can log the error or return it, depending on your application's needs
		log.Printf("Error adding account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add acccount"})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *Controller) DeleteAccountController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := models.DeleteAccount(id, userID.(int)) // Add account to storage
	if err != nil {
		// You can log the error or return it, depending on your application's needs
		log.Printf("Error adding account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete acccount"})
		return
	}
	c.JSON(http.StatusOK, account)
}
