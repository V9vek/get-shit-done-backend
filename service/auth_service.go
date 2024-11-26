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

	// generate refresh token
	refreshToken, err := s.JWTAuth.GenerateRefreshToken(fmt.Sprintf("%d", userId))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// generate access token
	accessToken, err := s.JWTAuth.GenerateAccessToken(fmt.Sprintf("%d", userId))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &model.Token{Refresh: refreshToken, Access: accessToken}, nil
}

func (s *AuthService) SignIn(context context.Context, username, password string) (*model.Token, error) {
	userId, err := s.AuthRepository.SignIn(context, username, password)
	if err != nil {
		return nil, fmt.Errorf("signin: %w", err)
	}

	// generate refresh token
	refreshToken, err := s.JWTAuth.GenerateRefreshToken(fmt.Sprintf("%d", userId))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// generate access token
	accessToken, err := s.JWTAuth.GenerateAccessToken(fmt.Sprintf("%d", userId))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &model.Token{Refresh: refreshToken, Access: accessToken}, nil
}

func (s *AuthService) DeleteRefreshToken(context context.Context, userId int, refreshToken string) error {
	isUserIdExist, err := s.AuthRepository.DoesUserIdExist(context, userId)
	if err != nil {
		return fmt.Errorf("user does not exist which has this refresh token: %w", err)
	}

	if isUserIdExist {
		err := s.AuthRepository.DeleteRefreshToken(userId, refreshToken)
		if err != nil {
			return fmt.Errorf("failed to delete refresh token: %w", err)
		}
		return nil
	}

	return fmt.Errorf("user does not exist")
}
