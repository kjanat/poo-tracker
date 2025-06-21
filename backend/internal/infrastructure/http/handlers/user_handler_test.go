package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// stubUserService implements user.Service with no-op methods for testing
// Only methods relevant to these tests need basic behavior

type stubUserService struct{}

func (s *stubUserService) Create(ctx context.Context, input *user.CreateUserInput) (*user.User, error) {
	return nil, nil
}
func (s *stubUserService) GetByID(ctx context.Context, id string) (*user.User, error) {
	return &user.User{ID: id}, nil
}
func (s *stubUserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
}
func (s *stubUserService) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return nil, nil
}
func (s *stubUserService) Update(ctx context.Context, id string, input *user.UpdateUserInput) (*user.User, error) {
	return &user.User{ID: id}, nil
}
func (s *stubUserService) Delete(ctx context.Context, id string) error { return nil }
func (s *stubUserService) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	return nil, nil
}
func (s *stubUserService) ListWithCount(ctx context.Context, limit, offset int) ([]*user.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserService) Register(ctx context.Context, input *user.RegisterInput) (*user.User, error) {
	return nil, nil
}
func (s *stubUserService) Login(ctx context.Context, input *user.LoginInput) (*user.LoginResult, error) {
	return nil, nil
}
func (s *stubUserService) ChangePassword(ctx context.Context, userID string, input *user.ChangePasswordInput) error {
	return nil
}
func (s *stubUserService) DeactivateAccount(ctx context.Context, userID string) error { return nil }
func (s *stubUserService) GetSettings(ctx context.Context, userID string) (*user.UserSettings, error) {
	return &user.UserSettings{UserID: userID}, nil
}
func (s *stubUserService) UpdateSettings(ctx context.Context, userID string, input *user.UpdateSettingsInput) (*user.UserSettings, error) {
	return &user.UserSettings{UserID: userID}, nil
}
func (s *stubUserService) ValidateEmail(ctx context.Context, email string) error       { return nil }
func (s *stubUserService) ValidateUsername(ctx context.Context, username string) error { return nil }
func (s *stubUserService) ValidatePassword(password string) error                      { return nil }

var _ user.Service = (*stubUserService)(nil)

func TestHandlers_InvalidUserIDContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubUserService{}
	h := NewUserHandler(svc)

	tests := []struct {
		name    string
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GetProfile", http.MethodGet, "/profile", h.GetProfile},
		{"UpdateProfile", http.MethodPut, "/profile", h.UpdateProfile},
		{"ChangePassword", http.MethodPost, "/change-password", h.ChangePassword},
		{"GetSettings", http.MethodGet, "/settings", h.GetSettings},
		{"UpdateSettings", http.MethodPut, "/settings", h.UpdateSettings},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user_id", 123) // invalid type
			})
			router.Handle(tt.method, tt.path, tt.handler)

			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected status 401, got %d", w.Code)
			}
		})
	}
}
