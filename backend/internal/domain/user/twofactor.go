package user

import "time"

// UserTwoFactor represents a user's two-factor authentication settings
// stored by the two-factor service.
type UserTwoFactor struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Secret      string     `json:"-"`
	IsEnabled   bool       `json:"isEnabled"`
	BackupCodes []string   `json:"-"`
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// UserTwoFactorSetupRequest represents a request to set up 2FA
// during the onboarding flow.
type UserTwoFactorSetupRequest struct {
	Token string `json:"token" binding:"required"`
}

// UserTwoFactorVerifyRequest represents a request to verify 2FA
// using either a TOTP token or a backup code.
type UserTwoFactorVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// UserTwoFactorSetupResponse represents the response for 2FA setup
// including the TOTP secret and backup codes.
type UserTwoFactorSetupResponse struct {
	Secret      string   `json:"secret"`
	BackupCodes []string `json:"backupCodes"`
	QRCodeURL   string   `json:"qrCodeUrl"`
}

// UserTwoFactorStatus represents the current 2FA status
// for the authenticated user.
type UserTwoFactorStatus struct {
	IsEnabled        bool       `json:"isEnabled"`
	LastUsedAt       *time.Time `json:"lastUsedAt,omitempty"`
	BackupCodesCount int        `json:"backupCodesCount"`
}
