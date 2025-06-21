package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func setupTwoFactorService() (*TwoFactorService, repository.UserTwoFactorRepository, *repository.MemoryUserRepository, *user.User) {
	userRepo := repository.NewMemoryUserRepository()
	twoRepo := repository.NewMemoryUserTwoFactorRepository()
	svc := NewTwoFactorService(twoRepo, userRepo, "PooTracker")
	u := &user.User{ID: "u1", Email: "user@example.com", Name: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_ = userRepo.CreateUser(u)
	return svc, twoRepo, userRepo, u
}

func TestTwoFactorService_GenerateSecret(t *testing.T) {
	svc, _, _, usr := setupTwoFactorService()
	ctx := context.Background()

	resp, err := svc.GenerateSecret(ctx, usr.ID)
	if err != nil {
		t.Fatalf("GenerateSecret failed: %v", err)
	}
	if resp.Secret == "" {
		t.Error("expected secret to be generated")
	}
	if len(resp.BackupCodes) != 10 {
		t.Errorf("expected 10 backup codes, got %d", len(resp.BackupCodes))
	}
	if !strings.Contains(resp.QRCodeURL, usr.Email) {
		t.Errorf("qr code url should contain email, got %s", resp.QRCodeURL)
	}
}

func TestTwoFactorService_EnableTwoFactor(t *testing.T) {
	svc, repo, _, usr := setupTwoFactorService()
	ctx := context.Background()

	resp, err := svc.GenerateSecret(ctx, usr.ID)
	if err != nil {
		t.Fatalf("GenerateSecret failed: %v", err)
	}
	code, err := totp.GenerateCode(resp.Secret, time.Now())
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if err := svc.EnableTwoFactor(ctx, usr.ID, code, resp.Secret, resp.BackupCodes); err != nil {
		t.Fatalf("EnableTwoFactor failed: %v", err)
	}

	tf, err := repo.GetByUserID(ctx, usr.ID)
	if err != nil {
		t.Fatalf("repo GetByUserID failed: %v", err)
	}
	if !tf.IsEnabled {
		t.Error("two factor should be enabled")
	}
	if tf.Secret != resp.Secret {
		t.Error("stored secret mismatch")
	}
	if len(tf.BackupCodes) != len(resp.BackupCodes) {
		t.Error("backup codes not stored correctly")
	}
}

func TestTwoFactorService_VerifyToken(t *testing.T) {
	svc, repo, _, usr := setupTwoFactorService()
	ctx := context.Background()

	resp, _ := svc.GenerateSecret(ctx, usr.ID)
	code, _ := totp.GenerateCode(resp.Secret, time.Now())
	_ = svc.EnableTwoFactor(ctx, usr.ID, code, resp.Secret, resp.BackupCodes)

	t.Run("valid TOTP", func(t *testing.T) {
		tok, _ := totp.GenerateCode(resp.Secret, time.Now())
		ok, err := svc.VerifyToken(ctx, usr.ID, tok)
		if err != nil {
			t.Fatalf("VerifyToken failed: %v", err)
		}
		if !ok {
			t.Error("expected token to be valid")
		}
		tf, _ := repo.GetByUserID(ctx, usr.ID)
		if tf.LastUsedAt == nil {
			t.Error("LastUsedAt should be set")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		ok, err := svc.VerifyToken(ctx, usr.ID, "000000")
		if err != nil {
			t.Fatalf("VerifyToken unexpected error: %v", err)
		}
		if ok {
			t.Error("expected invalid token to fail")
		}
	})

	t.Run("backup code", func(t *testing.T) {
		code := resp.BackupCodes[0]
		ok, err := svc.VerifyToken(ctx, usr.ID, code)
		if err != nil {
			t.Fatalf("VerifyToken failed: %v", err)
		}
		if !ok {
			t.Error("expected backup code to succeed")
		}
		tf, _ := repo.GetByUserID(ctx, usr.ID)
		if len(tf.BackupCodes) != len(resp.BackupCodes)-1 {
			t.Error("backup code should be removed after use")
		}
	})
}

func TestTwoFactorService_DisableTwoFactor(t *testing.T) {
	svc, repo, _, usr := setupTwoFactorService()
	ctx := context.Background()

	resp, _ := svc.GenerateSecret(ctx, usr.ID)
	code, _ := totp.GenerateCode(resp.Secret, time.Now())
	_ = svc.EnableTwoFactor(ctx, usr.ID, code, resp.Secret, resp.BackupCodes)

	if err := svc.DisableTwoFactor(ctx, usr.ID); err != nil {
		t.Fatalf("DisableTwoFactor failed: %v", err)
	}
	_, err := repo.GetByUserID(ctx, usr.ID)
	if err != repository.ErrNotFound {
		t.Errorf("expected ErrNotFound after disable, got %v", err)
	}
}

func TestTwoFactorService_GetStatus(t *testing.T) {
	svc, repo, _, usr := setupTwoFactorService()
	ctx := context.Background()

	status, err := svc.GetStatus(ctx, usr.ID)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.IsEnabled {
		t.Error("expected 2FA disabled by default")
	}

	resp, _ := svc.GenerateSecret(ctx, usr.ID)
	code, _ := totp.GenerateCode(resp.Secret, time.Now())
	_ = svc.EnableTwoFactor(ctx, usr.ID, code, resp.Secret, resp.BackupCodes)
	code, _ = totp.GenerateCode(resp.Secret, time.Now())
	_, _ = svc.VerifyToken(ctx, usr.ID, code)

	status, err = svc.GetStatus(ctx, usr.ID)
	if err != nil {
		t.Fatalf("GetStatus after enable failed: %v", err)
	}
	if !status.IsEnabled {
		t.Error("expected enabled status")
	}
	if status.BackupCodesCount != len(resp.BackupCodes) {
		t.Errorf("expected %d backup codes, got %d", len(resp.BackupCodes), status.BackupCodesCount)
	}
	if status.LastUsedAt == nil {
		t.Error("expected LastUsedAt to be set")
	}

	// Ensure repo LastUsedAt matches
	tf, _ := repo.GetByUserID(ctx, usr.ID)
	if tf.LastUsedAt == nil || status.LastUsedAt == nil {
		t.Fatal("missing LastUsedAt")
	}
	if !tf.LastUsedAt.Equal(*status.LastUsedAt) {
		t.Error("LastUsedAt mismatch")
	}
}
