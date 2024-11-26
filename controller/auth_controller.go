package controller

import (
	"fmt"
	"get-shit-done/model"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"strconv"
	"time"
)

type AuthController struct {
	authService          *service.AuthService
	JwtService           *service.JWTAuth
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthController(
	authService *service.AuthService,
	jwtService *service.JWTAuth,
	accessTokenDuration, refreshTokenDuration time.Duration,
) *AuthController {
	return &AuthController{
		authService:          authService,
		JwtService:           jwtService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (c *AuthController) SignUp(writer http.ResponseWriter, requests *http.Request) {
	var cred model.AuthCredentials
	if err := utils.ReadFromRequestBody(requests, &cred); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
		return
	}

	token, err := c.authService.SignUp(requests.Context(), cred)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to sign up: %v", err), http.StatusUnauthorized)
		return
	}

	c.setTokensInCookies(writer, token)

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully signed up",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *AuthController) SignIn(writer http.ResponseWriter, requests *http.Request) {
	var cred model.SigninCredentials
	if err := utils.ReadFromRequestBody(requests, &cred); err != nil {
		http.Error(writer, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
		return
	}

	token, err := c.authService.SignIn(requests.Context(), cred.Username, cred.Password)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to sign in: %v", err), http.StatusUnauthorized)
		return
	}

	c.setTokensInCookies(writer, token)
	fmt.Println(token.Access)

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully signed in",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *AuthController) setTokensInCookies(writer http.ResponseWriter, token *model.Token) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "access_token",
		Value:    token.Access,
		Path:     "/",
		Domain:   "",                                    // Set to your domain if needed
		Expires:  time.Now().Add(c.accessTokenDuration), // Set expiration as per your requirements
		Secure:   true,                                  // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.Refresh,
		Path:     "/",
		Domain:   "",                                     // Set to your domain if needed
		Expires:  time.Now().Add(c.refreshTokenDuration), // Set expiration as per your requirements
		Secure:   true,                                   // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AuthController) RefreshRefreshToken(writer http.ResponseWriter, requests *http.Request) {
	refreshTokenCookie, err := requests.Cookie("refresh_token")

	if err != nil {
		http.Error(writer, fmt.Sprintf("refresh token not found: %v", err), http.StatusUnauthorized)
		return
	}

	refreshToken := refreshTokenCookie.Value
	updatedToken, err := c.JwtService.RefreshRefreshToken(requests.Context(), refreshToken)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to refresh refresh token: %v", err), http.StatusUnauthorized)
		return
	}

	c.setTokensInCookies(writer, updatedToken)

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully refreshed the tokens",
		Data:   nil,
	}
	utils.WriteResponseBody(writer, webResponse)
}

func (c *AuthController) Logout(writer http.ResponseWriter, requests *http.Request) {
	// delete the refresh token
	refreshTokenCookie, err := requests.Cookie("refresh_token")

	if err != nil {
		http.Error(writer, "refresh token not found", http.StatusUnauthorized)
		return
	}

	refreshToken := refreshTokenCookie.Value

	sub, err := c.JwtService.GetSubjectFromRefreshToken(refreshToken)
	if err != nil {
		http.Error(writer, "failed to get the subject from token", http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		http.Error(writer, "token's subject has invalid format", http.StatusUnauthorized)
		return
	}

	err = c.authService.DeleteRefreshToken(requests.Context(), userId, refreshToken)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to log out: %v", err), http.StatusBadRequest)
		return
	}

	// Clear cookies
	http.SetCookie(writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   "",         // Set to your domain if needed
		Expires:  time.Now(), // Set expiration as per your requirements
		Secure:   true,       // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "",         // Set to your domain if needed
		Expires:  time.Now(), // Set expiration as per your requirements
		Secure:   true,       // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c *AuthController) VerifyAuth(writer http.ResponseWriter, requests *http.Request) {
	accessTokenCookie, err := requests.Cookie("access_token")
	if err != nil {
		http.Error(writer, fmt.Sprintf("access token not found: %v", err), http.StatusUnauthorized)
		return
	}

	accessToken := accessTokenCookie.Value

	isValid, err := c.JwtService.IsAccessTokenValid(accessToken)
	if !isValid || err != nil {
		http.Error(writer, fmt.Sprintf("invalid or expired token %v", err), http.StatusUnauthorized)
		return
	}
}
