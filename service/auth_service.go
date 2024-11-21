package service

import (
	"context"
	"fmt"
	"get-shit-done/model"
	"get-shit-done/repository"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
	JWTAuth        *JWTAuth
}

func NewAuthService(authRepository *repository.AuthRepository, jwtAuth *JWTAuth) *AuthService {
	return &AuthService{AuthRepository: authRepository, JWTAuth: jwtAuth}
}

func (s *AuthService) SignUp(context context.Context, credentials model.AuthCredentials) (*model.Token, error) {
	isUsernameExist, err := s.AuthRepository.DoesUsernameExist(context, credentials.Username)
	if err != nil {
		return nil, fmt.Errorf("signup: %w", err)
	}

	if isUsernameExist {
		return nil, fmt.Errorf("signup: username already exist")
	}

	userId, err := s.AuthRepository.SignUp(context, credentials)
	if err != nil {
		return nil, fmt.Errorf("signup failed: %w", err)
	}

	// generate refresh and access token
}

func (s *AuthService) SignIn(context context.Context, username, password string) (*model.Token, error) {
	userId, err := s.AuthRepository.SignIn(context, username, password)
	if err != nil {
		return nil, fmt.Errorf("signin: %w", err)
	}

	// generate refresh and access token
}
