package model

import (
	"time"
)

// UserTwoFactor represents a user's two-factor authentication settings
type UserTwoFactor struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Secret      string     `json:"-"` // Never expose the secret in JSON
	IsEnabled   bool       `json:"isEnabled"`
	BackupCodes []string   `json:"-"` // Never expose backup codes in JSON
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// UserTwoFactorSetupRequest represents a request to set up 2FA
type UserTwoFactorSetupRequest struct {
	Token string `json:"token" binding:"required"` // TOTP token for verification
}

// UserTwoFactorVerifyRequest represents a request to verify 2FA
type UserTwoFactorVerifyRequest struct {
	Token string `json:"token" binding:"required"` // TOTP token or backup code
}

// UserTwoFactorSetupResponse represents the response for 2FA setup
type UserTwoFactorSetupResponse struct {
	Secret      string   `json:"secret"`      // TOTP secret for QR code generation
	BackupCodes []string `json:"backupCodes"` // One-time backup codes
	QRCodeURL   string   `json:"qrCodeUrl"`   // URL for QR code generation
}

// UserTwoFactorStatus represents the current 2FA status
type UserTwoFactorStatus struct {
	IsEnabled        bool       `json:"isEnabled"`
	LastUsedAt       *time.Time `json:"lastUsedAt,omitempty"`
	BackupCodesCount int        `json:"backupCodesCount"` // Number of unused backup codes
}
