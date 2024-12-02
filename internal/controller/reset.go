package controller

import (
	"guilliman/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) ResetController(c *gin.Context) {
	err := models.Reset()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reset successful"})
}
