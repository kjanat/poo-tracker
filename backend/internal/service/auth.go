package service

import (
	"poo-tracker/internal/model"
)

type AuthService interface {
	Register(email, username, password string) (*model.User, error)
	Login(email, password string) (string, error) // returns JWT
	ValidateToken(token string) (*model.User, error)
}

// JWTAuthService is a stub for JWT-based authentication.
type JWTAuthService struct{}

func (a *JWTAuthService) Register(email, username, password string) (*model.User, error) {
	// TODO: implement registration logic
	return nil, nil
}

func (a *JWTAuthService) Login(email, password string) (string, error) {
	// TODO: implement login logic
	return "", nil
}

func (a *JWTAuthService) ValidateToken(token string) (*model.User, error) {
	// TODO: implement token validation
	return nil, nil
}
