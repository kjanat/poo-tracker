package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// UserRepository implements user.Repository using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, update *user.UserUpdate) error {
	updates := make(map[string]interface{})

	if update.Email != nil {
		updates["email"] = *update.Email
	}
	if update.Username != nil {
		updates["username"] = *update.Username
	}
	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.Age != nil {
		updates["age"] = *update.Age
	}
	if update.Gender != nil {
		updates["gender"] = *update.Gender
	}
	if update.Height != nil {
		updates["height"] = *update.Height
	}
	if update.Weight != nil {
		updates["weight"] = *update.Weight
	}

	return r.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	var users []*user.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// UserAuth operations
func (r *UserRepository) CreateAuth(ctx context.Context, auth *user.UserAuth) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

func (r *UserRepository) GetAuthByUserID(ctx context.Context, userID string) (*user.UserAuth, error) {
	var auth user.UserAuth
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&auth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserAuthNotFound
		}
		return nil, err
	}
	return &auth, nil
}

func (r *UserRepository) UpdateAuth(ctx context.Context, userID string, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&user.UserAuth{}).Where("user_id = ?", userID).Update("password_hash", passwordHash).Error
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&user.UserAuth{}).Where("user_id = ?", userID).Update("last_login_at", "NOW()").Error
}

func (r *UserRepository) DeactivateAuth(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&user.UserAuth{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

// UserSettings operations
func (r *UserRepository) CreateSettings(ctx context.Context, settings *user.UserSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

func (r *UserRepository) GetSettingsByUserID(ctx context.Context, userID string) (*user.UserSettings, error) {
	var settings user.UserSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserSettingsNotFound
		}
		return nil, err
	}
	return &settings, nil
}

func (r *UserRepository) UpdateSettings(ctx context.Context, userID string, update *user.UserSettingsUpdate) error {
	updates := make(map[string]interface{})

	if update.Timezone != nil {
		updates["timezone"] = *update.Timezone
	}
	if update.ReminderTime != nil {
		updates["reminder_time"] = *update.ReminderTime
	}
	if update.ReminderEnabled != nil {
		updates["reminder_enabled"] = *update.ReminderEnabled
	}
	if update.DataRetentionDays != nil {
		updates["data_retention_days"] = *update.DataRetentionDays
	}
	if update.PrivacyLevel != nil {
		updates["privacy_level"] = *update.PrivacyLevel
	}
	if update.NotificationEnabled != nil {
		updates["notification_enabled"] = *update.NotificationEnabled
	}
	if update.ThemePreference != nil {
		updates["theme_preference"] = *update.ThemePreference
	}
	if update.DarkMode != nil {
		updates["dark_mode"] = *update.DarkMode
	}
	if update.PreferredUnits != nil {
		updates["preferred_units"] = *update.PreferredUnits
	}

	return r.db.WithContext(ctx).Model(&user.UserSettings{}).Where("user_id = ?", userID).Updates(updates).Error
}

// Query operations
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error
	return count, err
}
