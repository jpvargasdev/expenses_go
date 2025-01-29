package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserUID(c *gin.Context) (string, error) {
	userUID, exists := c.Get("userUID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return "", fmt.Errorf("User not authenticated")
	}

	uid, ok := userUID.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user UID"})
		return "", fmt.Errorf("Failed to get user UID")
	}

	return uid, nil
}
