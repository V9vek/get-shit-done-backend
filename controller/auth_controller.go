package controller

import (
	"fmt"
	"get-shit-done/model"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"time"
)

type AuthController struct {
	authService          *service.AuthService
	jwtService           *service.JWTAuth
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
		jwtService:           jwtService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (c *AuthController) SignUp(writer http.ResponseWriter, requests *http.Request) {
	var cred model.AuthCredentials
	if err := utils.ReadFromRequestBody(requests, &cred); err != nil {
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("Invalid credentials: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
		return
	}

	token, err := c.authService.SignUp(requests.Context(), cred)
	if err != nil {
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("Failed to signup: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
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
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("Invalid credentials: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
		return
	}

	token, err := c.authService.SignIn(requests.Context(), cred.Username, cred.Password)
	if err != nil {
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("Failed to signin: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
		return
	}

	c.setTokensInCookies(writer, token)

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
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("refresh token not found: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
		return
	}

	refreshToken := refreshTokenCookie.Value
	updatedToken, err := c.jwtService.RefreshRefreshToken(requests.Context(), refreshToken)
	if err != nil {
		webResponse := model.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: fmt.Sprintf("failed to refresh refresh token: %v", err),
			Data:   nil,
		}
		utils.WriteResponseBody(writer, webResponse)
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
