package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

// TwoFactorService provides two-factor authentication functionality
type TwoFactorService struct {
	repo     repository.UserTwoFactorRepository
	userRepo repository.UserRepository
	issuer   string // The name of your application for TOTP
}

var ErrInvalidToken = errors.New("invalid token")

// NewTwoFactorService creates a new TwoFactorService
func NewTwoFactorService(repo repository.UserTwoFactorRepository, userRepo repository.UserRepository, issuer string) *TwoFactorService {
	return &TwoFactorService{
		repo:     repo,
		userRepo: userRepo,
		issuer:   issuer,
	}
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TwoFactorService) GenerateSecret(ctx context.Context, userID string) (*user.UserTwoFactorSetupResponse, error) {
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

	return &user.UserTwoFactorSetupResponse{
		Secret:      secret,
		BackupCodes: backupCodes,
		QRCodeURL:   qrCodeURL,
	}, nil
}

// EnableTwoFactor enables 2FA for a user after verifying the initial token
func (s *TwoFactorService) EnableTwoFactor(ctx context.Context, userID string, token string, secret string, backupCodes []string) error {
	// Verify the token before enabling
	if !s.verifyTOTPToken(secret, token) {
		return ErrInvalidToken
	}

	// Create or update the 2FA record
	twoFactor := &user.UserTwoFactor{
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
func (s *TwoFactorService) GetStatus(ctx context.Context, userID string) (*user.UserTwoFactorStatus, error) {
	twoFactor, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return &user.UserTwoFactorStatus{
				IsEnabled:        false,
				BackupCodesCount: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get 2FA settings: %w", err)
	}

	return &user.UserTwoFactorStatus{
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
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}
		code := base32.StdEncoding.EncodeToString(bytes)
		code = strings.ToUpper(code)[:8] // Take first 8 characters
		// Format as XXXX-XXXX for readability
		codes[i] = code[:4] + "-" + code[4:]
	}
	return codes, nil
}

// generateQRCodeURL generates a QR code URL for Google Authenticator
func (s *TwoFactorService) generateQRCodeURL(email string, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		s.issuer, email, secret, s.issuer)
}

// verifyTOTPToken verifies a TOTP token against a secret
func (s *TwoFactorService) verifyTOTPToken(secret string, token string) bool {
	// Use the pquerna/otp library's TOTP validation
	valid, err := totp.ValidateCustom(token, secret, time.Now(), totp.ValidateOpts{
		Period:    30, // Default period of 30 seconds
		Skew:      1,  // Allow 1 period skew for clock drift
		Digits:    6,  // Use 6 digit tokens
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}
	return valid
}
