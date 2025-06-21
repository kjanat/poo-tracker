package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth represents authentication data for a user
// UserAuth is used for password hash storage, not User struct.
type AuthService interface {
	Register(email, password, name string) (*user.User, string, error)
	Login(email, password string) (*user.User, string, error)
	ValidateToken(token string) (*user.User, error)
}

type JWTAuthService struct {
	UserRepo repository.UserRepository
	Secret   string
	Expiry   time.Duration
}

func (a *JWTAuthService) Register(email, password, name string) (*user.User, string, error) {
	// Check if email is already registered, distinguishing between not found and other errors
	_, err := a.UserRepo.GetUserByEmail(email)
	if err == nil {
		return nil, "", errors.New("email already registered")
	}
	// Only accept "not found" error - return other errors directly
	if err != repository.ErrNotFound {
		return nil, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user := &user.User{
		ID:        uuid.NewString(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	err = a.UserRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	// Create auth data - if this fails, we should clean up the user
	auth := &user.UserAuth{
		UserID:       user.ID,
		PasswordHash: string(hash),
		Provider:     "local",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = a.UserRepo.CreateUserAuth(auth)
	if err != nil {
		// Clean up the user if auth creation fails
		_ = a.UserRepo.DeleteUser(user.ID)
		return nil, "", err
	}

	token, err := a.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (a *JWTAuthService) Login(email, password string) (*user.User, string, error) {
	user, err := a.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	auth, err := a.UserRepo.GetUserAuthByUserID(user.ID)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := a.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (a *JWTAuthService) ValidateToken(tokenStr string) (*user.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm to prevent algorithm confusion attacks
		if token.Method.Alg() != "HS256" {
			return nil, errors.New("invalid signing algorithm")
		}
		return []byte(a.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token subject")
	}
	return a.UserRepo.GetUserByID(userID)
}

func (a *JWTAuthService) generateToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(a.Expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.Secret))
}
