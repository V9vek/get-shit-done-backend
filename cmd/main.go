package main

import (
	"fmt"
	"get-shit-done/config"
	"get-shit-done/routes"
	"get-shit-done/utils"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	// db connection
	db := config.DatabaseConnection()
	defer db.Close()

	// router setup
	r := routes.SetupRoutes()

	// server setup
	port := ":8080"
	err := http.ListenAndServe(port, r)
	utils.PanicIfError(err)
}
