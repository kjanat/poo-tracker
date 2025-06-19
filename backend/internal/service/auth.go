package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth represents authentication data for a user
// UserAuth is used for password hash storage, not User struct.
type AuthService interface {
	Register(email, password, name string) (*model.User, string, error)
	Login(email, password string) (*model.User, string, error)
	ValidateToken(token string) (*model.User, error)
}

type JWTAuthService struct {
	UserRepo repository.UserRepository
	Secret   string
	Expiry   time.Duration
}

func (a *JWTAuthService) Register(email, password, name string) (*model.User, string, error) {
	if _, err := a.UserRepo.GetUserByEmail(email); err == nil {
		return nil, "", errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user := &model.User{
		ID:        uuid.NewString(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = a.UserRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}
	auth := &model.UserAuth{
		UserID:       user.ID,
		PasswordHash: string(hash),
		Provider:     "local",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = a.UserRepo.CreateUserAuth(auth)
	if err != nil {
		return nil, "", err
	}
	token, err := a.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (a *JWTAuthService) Login(email, password string) (*model.User, string, error) {
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

func (a *JWTAuthService) ValidateToken(tokenStr string) (*model.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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

func (a *JWTAuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(a.Expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.Secret))
}
