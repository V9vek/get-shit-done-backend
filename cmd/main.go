package main

import (
	"fmt"
	"get-shit-done/config"
	"get-shit-done/controller"
	"get-shit-done/repository"
	"get-shit-done/routes"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
	"os"
	"time"
)

var (
	jwtSecretKeyAccess     = os.Getenv("JWT_SECRET_KEY_ACCESS")
	jwtSecretKeyRefresh    = os.Getenv("JWT_SECRET_KEY_REFRESH")
	jwtIss                 = os.Getenv("JWT_ISS")
	jwtAccessTokenExpTime  = os.Getenv("JWT_ACCESS_TOKEN_EXP_TIME")
	jwtRefreshTokenExpTime = os.Getenv("JWT_REFRESH_TOKEN_EXP_TIME")
)

func main() {
	fmt.Println("Starting server...")

	// db
	db := config.DatabaseConnection()
	defer db.Close()

	// repository
	authRepository := repository.NewAuthRepository(db)

	// jwt
	jwtAccessTokenExpTimeDuration, err := time.ParseDuration("5m")
	utils.PanicIfError(err)
	jwtRefreshTokenExpTimeDuration, err := time.ParseDuration("24h")
	utils.PanicIfError(err)

	jwtAuth := service.NewJWTAuth(
		*authRepository,
		jwtSecretKeyAccess,
		jwtSecretKeyRefresh,
		jwtIss,
		jwtAccessTokenExpTimeDuration,
		jwtRefreshTokenExpTimeDuration,
	)

	// service
	authService := service.NewAuthService(authRepository, jwtAuth)

	// controller
	authController := controller.NewAuthController(
		authService,
		jwtAuth,
		jwtAccessTokenExpTimeDuration,
		jwtRefreshTokenExpTimeDuration,
	)

	// router
	r := routes.SetupRoutes(authController)

	// server
	port := ":8080"
	err = http.ListenAndServe(port, r)
	utils.PanicIfError(err)
}
