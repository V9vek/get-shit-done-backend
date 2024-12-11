package service

import (
	"context"
	"errors"
	"fmt"
	"get-shit-done/model"
	"get-shit-done/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTAuth struct {
	authRepo            repository.AuthRepository
	secretKeyAccess     string
	secretKeyRefresh    string
	iss                 string
	accessTokenExpTime  time.Duration
	refreshTokenExpTime time.Duration
}

func NewJWTAuth(
	authRepo repository.AuthRepository,
	secretKeyAcess, secretKeyRefresh, iss string,
	accessTokenExpTime, refreshTokenExpTime time.Duration,
) *JWTAuth {
	return &JWTAuth{
		authRepo:            authRepo,
		secretKeyAccess:     secretKeyAcess,
		secretKeyRefresh:    secretKeyRefresh,
		iss:                 iss,
		accessTokenExpTime:  accessTokenExpTime,
		refreshTokenExpTime: refreshTokenExpTime,
	}
}

func parse(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j JWTAuth) IsAccessTokenValid(accessToken string) (bool, error) {
	token, err := parse(accessToken, j.secretKeyAccess)
	if err != nil {
		return false, fmt.Errorf("can not parse the access token %w", err)
	}

	isValid, _ := isValidTime(token)

	if !isValid {
		return false, fmt.Errorf("access token is expired")
	}
	return true, nil
}

func (j JWTAuth) IsRefreshTokenValid(context context.Context, refreshToken string) (bool, error) {
	token, err := parse(refreshToken, j.secretKeyAccess)
	if err != nil {
		return false, fmt.Errorf("can not parse the refresh token %w", err)
	}

	isValid, err := isValidTime(token)

	if !isValid {
		return false, fmt.Errorf("refresh token is expired: %w", err)
	}

	// check if the refresh token is of the requested user or not
	sub, err := getSubjectFromToken(token)
	if err != nil {
		return false, fmt.Errorf("failed to get the subject from token: %w", err)
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		return false, fmt.Errorf("token's subject has invalid format: %w", err)
	}

	isUserIdExist, err := j.authRepo.DoesUserIdExist(context, userId)
	if err != nil {
		return false, fmt.Errorf("user does not exist which has this refresh token: %w", err)
	}

	return isUserIdExist, nil
}

func (j JWTAuth) GetSubjectFromAccessToken(tokenStr string) (string, error) {
	token, err := parse(tokenStr, j.secretKeyAccess)
	if err != nil {
		return "", fmt.Errorf("couldn't parse token string and secret: %w", err)
	}
	return getSubjectFromToken(token)
}

func (j JWTAuth) GetSubjectFromRefreshToken(tokenStr string) (string, error) {
	token, err := parse(tokenStr, j.secretKeyRefresh)
	if err != nil {
		return "", fmt.Errorf("couldn't parse token string and secret: %w", err)
	}
	return getSubjectFromToken(token)
}

func getSubjectFromToken(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", fmt.Errorf("`sub` claim is not present in token")
	}
	return "", fmt.Errorf("invalid token claims or token is not valid")
}

func isValidTime(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("exp property has invalid value")
		}
		now := float64(time.Now().Unix())
		if exp > now {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("invalid token claims or token is not valid")
}

// refreshing the Refresh Token = New Refresh Token + New Access Token also
func (j JWTAuth) RefreshRefreshToken(context context.Context, refreshToken string) (*model.Token, error) {
	isValid, err := j.IsRefreshTokenValid(context, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	// if refresh token expires and someone tries to refresh it, then it means security error
	// again the user have to sign up
	if !isValid {
		return nil, errors.New("invalid refresh token")
	}

	token, err := parse(refreshToken, j.secretKeyRefresh)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse token string and secret: %w", err)
	}

	sub, err := getSubjectFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get sub from refresh token: %w", err)
	}

	updatedRefreshTokenStr, err := j.GenerateRefreshToken(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	updatedAccessTokenStr, err := j.GenerateAccessToken(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// update the refresh token in DB
	userId, err := strconv.Atoi(sub)
	if err != nil {
		return nil, fmt.Errorf("sub has unsupported format: %w", err)
	}

	err = j.authRepo.UpdateRefreshToken(userId, updatedRefreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token in db: %w", err)
	}

	return &model.Token{Refresh: updatedRefreshTokenStr, Access: updatedAccessTokenStr}, nil
}

func (j JWTAuth) GenerateAccessToken(subject string) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    j.iss,
		Subject:   subject,
		ExpiresAt: time.Now().Add(j.accessTokenExpTime).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(j.secretKeyAccess))
	if err != nil {
		return "", fmt.Errorf("could not generate access token")
	}
	return token, nil
}

func (j JWTAuth) GenerateRefreshToken(subject string) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    j.iss,
		Subject:   subject,
		ExpiresAt: time.Now().Add(j.refreshTokenExpTime).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(j.secretKeyRefresh))
	if err != nil {
		return "", fmt.Errorf("could not generate refresh token")
	}

	// store/update the refresh token in DB
	userId, err := strconv.Atoi(subject)
	if err != nil {
		return "", fmt.Errorf("sub has unsupported format: %w", err)
	}

	err = j.authRepo.UpdateRefreshToken(userId, token)
	if err != nil {
		return "", fmt.Errorf("failed to update refresh token in db: %w", err)
	}

	return token, nil
}
