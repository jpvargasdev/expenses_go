package middleware

import (
	"context"
	"guilliman/cmd/firebase"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
    }

    // Extract the token
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Get Firebase Auth client
		client, err := firebase.FirebaseApp.Auth(context.Background())
		if err != nil {
      log.Printf("Error getting Firebase Auth: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase Auth"})
			c.Abort()
			return
		}

		// Verify the token
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Attach the UID to the context
		c.Set("userUID", token.UID)
		c.Next()
  }
}
