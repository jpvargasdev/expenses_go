package controller

import (
	"guilliman/internal/models"
	"guilliman/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) RegisterUserController(c *gin.Context) {
  var user models.User
  if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  user, err := models.RegisterUser(user)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
  }

  c.JSON(http.StatusCreated, gin.H{"status": "User created"})
}

func (h *Controller) LoginUserController(c *gin.Context) {
  var user models.User
  var userLogged models.UserLogged

  if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  userLogged, err := models.LoginUser(user)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in user"})
  }

  token, err := utils.GenerateToken(userLogged.Id)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging user"})
  }

  c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Controller) ResetUserController(c *gin.Context) {

}
