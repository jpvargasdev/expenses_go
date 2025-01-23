package controller 

import (
	"guilliman/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) DeleteUserController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
  if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }

  err := models.DeleteUser(uid)
  if err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to elete user"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (h *Controller) CreateUserController(c *gin.Context) {
  userUID, exists := c.Get("userUID")
  if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
  uid, ok := userUID.(string)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
    return
  }

  var newUser models.User
  if err := c.ShouldBindJSON(&newUser); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
  }

  newUser.ID = uid 

  err := models.CreateUser(newUser)
  if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, "OK")
}
