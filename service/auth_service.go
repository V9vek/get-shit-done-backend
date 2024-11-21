package service

import "get-shit-done/repository"

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{AuthRepository: authRepository}
}
