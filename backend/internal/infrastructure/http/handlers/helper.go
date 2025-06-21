package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// extractUserIDFromContext retrieves the user_id from context as a string.
// Returns the user ID and an error if missing or invalid type.
func extractUserIDFromContext(c *gin.Context) (string, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return "", errors.New("user_id not found in context")
    }
    userIDStr, ok := userID.(string)
    if !ok {
        return "", errors.New("user_id is not a string")
    }
    return userIDStr, nil
}

// extractUserID retrieves the user_id from context and handles HTTP errors.
func extractUserID(c *gin.Context) (string, bool) {
    userIDStr, err := extractUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return "", false
    }
    return userIDStr, true
}
