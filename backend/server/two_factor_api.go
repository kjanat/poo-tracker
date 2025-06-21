package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

// TwoFactorHandler handles two-factor authentication endpoints
type TwoFactorHandler struct {
	service *service.TwoFactorService
}

// NewTwoFactorHandler creates a new TwoFactorHandler
func NewTwoFactorHandler(service *service.TwoFactorService) *TwoFactorHandler {
	return &TwoFactorHandler{
		service: service,
	}
}

// GetStatus returns the 2FA status for the authenticated user
// GET /api/2fa/status
func (h *TwoFactorHandler) GetStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	status, err := h.service.GetStatus(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get 2FA status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// Setup generates a new 2FA secret and backup codes for setup
// POST /api/2fa/setup
func (h *TwoFactorHandler) Setup(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	setup, err := h.service.GenerateSecret(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA setup"})
		return
	}

	c.JSON(http.StatusOK, setup)
}

// Enable enables 2FA for the authenticated user
// POST /api/2fa/enable
func (h *TwoFactorHandler) Enable(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Secret      string   `json:"secret" binding:"required"`
		Token       string   `json:"token" binding:"required"`
		BackupCodes []string `json:"backupCodes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.EnableTwoFactor(c.Request.Context(), userID.(string), req.Token, req.Secret, req.BackupCodes)
	if err != nil {
		if err.Error() == "invalid token" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable disables 2FA for the authenticated user
// POST /api/2fa/disable
func (h *TwoFactorHandler) Disable(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req user.UserTwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Verify the token before disabling
	valid, err := h.service.VerifyToken(c.Request.Context(), userID.(string), req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	err = h.service.DisableTwoFactor(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Verify verifies a 2FA token for the authenticated user
// POST /api/2fa/verify
func (h *TwoFactorHandler) Verify(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req user.UserTwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	valid, err := h.service.VerifyToken(c.Request.Context(), userID.(string), req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": valid,
	})
}
