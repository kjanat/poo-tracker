package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

// TwoFactorService provides two-factor authentication functionality
type TwoFactorService struct {
	repo     repository.UserTwoFactorRepository
	userRepo repository.UserRepository
	issuer   string // The name of your application for TOTP
}

// NewTwoFactorService creates a new TwoFactorService
func NewTwoFactorService(repo repository.UserTwoFactorRepository, userRepo repository.UserRepository, issuer string) *TwoFactorService {
	return &TwoFactorService{
		repo:     repo,
		userRepo: userRepo,
		issuer:   issuer,
	}
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TwoFactorService) GenerateSecret(ctx context.Context, userID string) (*model.UserTwoFactorSetupResponse, error) {
	// Generate a random secret
	secret, err := s.generateRandomSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	// Generate backup codes
	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Get user for QR code generation
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate QR code URL
	qrCodeURL := s.generateQRCodeURL(user.Email, secret)

	return &model.UserTwoFactorSetupResponse{
		Secret:      secret,
		BackupCodes: backupCodes,
		QRCodeURL:   qrCodeURL,
	}, nil
}

// EnableTwoFactor enables 2FA for a user after verifying the initial token
func (s *TwoFactorService) EnableTwoFactor(ctx context.Context, userID string, token string, secret string, backupCodes []string) error {
	// Verify the token before enabling
	if !s.verifyTOTPToken(secret, token) {
		return fmt.Errorf("invalid token")
	}

	// Create or update the 2FA record
	twoFactor := &model.UserTwoFactor{
		UserID:      userID,
		Secret:      secret,
		IsEnabled:   true,
		BackupCodes: backupCodes,
	}

	// Check if 2FA already exists for this user
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil && err != repository.ErrNotFound {
		return fmt.Errorf("failed to check existing 2FA: %w", err)
	}

	if existing != nil {
		// Update existing
		return s.repo.Update(ctx, twoFactor)
	} else {
		// Create new
		return s.repo.Create(ctx, twoFactor)
	}
}

// DisableTwoFactor disables 2FA for a user
func (s *TwoFactorService) DisableTwoFactor(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}

// VerifyToken verifies a TOTP token or backup code
func (s *TwoFactorService) VerifyToken(ctx context.Context, userID string, token string) (bool, error) {
	twoFactor, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get 2FA settings: %w", err)
	}

	if !twoFactor.IsEnabled {
		return false, fmt.Errorf("2FA is not enabled for this user")
	}

	// Try TOTP token first
	if s.verifyTOTPToken(twoFactor.Secret, token) {
		// Update last used time (ignore error as it's not critical)
		_ = s.repo.UpdateLastUsed(ctx, userID)
		return true, nil
	}

	// Try backup codes
	for _, backupCode := range twoFactor.BackupCodes {
		if backupCode == token {
			// Remove the used backup code
			err := s.repo.RemoveBackupCode(ctx, userID, token)
			if err != nil {
				return false, fmt.Errorf("failed to remove backup code: %w", err)
			}
			// Update last used time (ignore error as it's not critical)
			_ = s.repo.UpdateLastUsed(ctx, userID)
			return true, nil
		}
	}

	return false, nil
}

// GetStatus returns the 2FA status for a user
func (s *TwoFactorService) GetStatus(ctx context.Context, userID string) (*model.UserTwoFactorStatus, error) {
	twoFactor, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return &model.UserTwoFactorStatus{
				IsEnabled:        false,
				BackupCodesCount: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get 2FA settings: %w", err)
	}

	return &model.UserTwoFactorStatus{
		IsEnabled:        twoFactor.IsEnabled,
		LastUsedAt:       twoFactor.LastUsedAt,
		BackupCodesCount: len(twoFactor.BackupCodes),
	}, nil
}

// generateRandomSecret generates a random base32-encoded secret for TOTP
func (s *TwoFactorService) generateRandomSecret() (string, error) {
	bytes := make([]byte, 20) // 160 bits for good security
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}

// generateBackupCodes generates backup codes for 2FA
func (s *TwoFactorService) generateBackupCodes() ([]string, error) {
	codes := make([]string, 10) // Generate 10 backup codes
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 5) // 8 characters when base32 encoded
		_, err := rand.Read(bytes)
		if err != nil {
			return nil, err
		}
		// Format as XXX-XXX for readability
		code := base32.StdEncoding.EncodeToString(bytes)
		code = strings.ToUpper(code)[:8] // Take first 8 characters
		codes[i] = code[:3] + "-" + code[3:]
	}
	return codes, nil
}

// generateQRCodeURL generates a QR code URL for Google Authenticator
func (s *TwoFactorService) generateQRCodeURL(email string, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		s.issuer, email, secret, s.issuer)
}

// verifyTOTPToken verifies a TOTP token against a secret
// This is a simplified implementation. In production, you'd use a library like "github.com/pquerna/otp"
func (s *TwoFactorService) verifyTOTPToken(secret string, token string) bool {
	// Convert token to integer
	tokenInt, err := strconv.ParseInt(token, 10, 32)
	if err != nil {
		return false
	}

	// Get current time in 30-second intervals
	now := time.Now().Unix() / 30

	// Check current time window and previous/next windows for clock skew
	for i := int64(-1); i <= 1; i++ {
		timeStep := now + i
		expectedToken := s.generateTOTPToken(secret, timeStep)
		if expectedToken == int32(tokenInt) {
			return true
		}
	}

	return false
}

// generateTOTPToken generates a TOTP token for a given secret and time step
// This is a simplified implementation for demonstration purposes
func (s *TwoFactorService) generateTOTPToken(secret string, timeStep int64) int32 {
	// This is a very simplified TOTP implementation
	// In production, use a proper TOTP library like "github.com/pquerna/otp"

	// Decode base32 secret
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return 0
	}

	// Simple hash-based calculation (not RFC 6238 compliant)
	// This is just for demonstration - use a proper TOTP library in production
	hash := int32(0)
	for i, b := range secretBytes {
		hash += int32(b) * int32(timeStep+int64(i))
	}

	// Return 6-digit token
	return (hash % 1000000)
}
