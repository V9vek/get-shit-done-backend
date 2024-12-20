package controller

import (
	"fmt"
	"get-shit-done/model"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"strconv"
	"strings"
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

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully signed up",
		Data:   token.Access,
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

	fmt.Println(token.Access)

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully signed in",
		Data:   token.Access,
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

	// remove this later
	c.setTokensInCookies(writer, updatedToken)

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully refreshed the tokens",
		Data:   updatedToken.Access,
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
		Secure:   false,      // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "",         // Set to your domain if needed
		Expires:  time.Now(), // Set expiration as per your requirements
		Secure:   false,      // Set to true if using HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getAccessTokenFromHeaders(w http.ResponseWriter, r *http.Request) string {
	auth := r.Header.Get("Authorization")

	if auth == "" {
		http.Error(w, "missing or malformed token", http.StatusUnauthorized)
		return ""
	}

	headerParts := strings.Split(auth, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(w, "missing or malformed token", http.StatusUnauthorized)
		return ""
	}

	token := headerParts[1]
	return token
}

func (c *AuthController) VerifyAuth(writer http.ResponseWriter, requests *http.Request) {
	auth := requests.Header.Get("Authorization")

	if auth == "" {
		http.Error(writer, "missing or malformed token", http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(auth, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(writer, "missing or malformed token", http.StatusUnauthorized)
		return
	}

	accessToken := headerParts[1]

	isValid, err := c.JwtService.IsAccessTokenValid(accessToken)
	if !isValid || err != nil {
		if err.Error() == "access token is expired" {
			sub, err := c.JwtService.GetSubjectFromAccessToken(accessToken)
			if err != nil {
				http.Error(writer, fmt.Sprintf("can't get sub from accessToken: %v", err), http.StatusUnauthorized)
				return
			}

			userId, err := strconv.Atoi(sub)
			if err != nil {
				http.Error(writer, fmt.Sprintf("not valid subject: %v", err), http.StatusUnauthorized)
				return
			}

			refreshToken, err := c.authService.GetRefreshToken(requests.Context(), userId)
			if err != nil {
				http.Error(writer, fmt.Sprintf("can't get refresh token from db: %v", err), http.StatusUnauthorized)
				return
			}

			tokens, err := c.JwtService.RefreshRefreshToken(requests.Context(), refreshToken)
			if err != nil {
				http.Error(writer, fmt.Sprintf("can't refresh the tokens: %v", err), http.StatusUnauthorized)
				return
			}

			webResponse := model.WebResponse{
				Code:   http.StatusOK,
				Status: "successfully refreshed the token",
				Data:   tokens.Access,
			}
			utils.WriteResponseBody(writer, webResponse)
			return
		} else {
			http.Error(writer, fmt.Sprintf("invalid token %v", err), http.StatusUnauthorized)
			return
		}
	}

	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "successfully validated the token",
		Data:   accessToken,
	}
	utils.WriteResponseBody(writer, webResponse)
}
