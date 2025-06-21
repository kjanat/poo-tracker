package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// extractUserID retrieves the user_id from the context ensuring it is a string.
// If the value is missing or not a string, it writes a 401 response and returns false.
func extractUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return "", false
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return "", false
	}
	return userIDStr, true
}
