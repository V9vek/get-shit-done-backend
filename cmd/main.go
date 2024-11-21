package main

import (
	"fmt"
	"get-shit-done/config"
	"get-shit-done/repository"
	"get-shit-done/routes"
	"get-shit-done/service"
	"get-shit-done/utils"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	// db
	db := config.DatabaseConnection()
	defer db.Close()

	// repository
	authRepository := repository.NewAuthRepository(db)

	// service
	authService := service.NewAuthService(authRepository)

	// router
	r := routes.SetupRoutes(authService)

	// server
	port := ":8080"
	err := http.ListenAndServe(port, r)
	utils.PanicIfError(err)
}
