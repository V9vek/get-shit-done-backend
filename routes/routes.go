package routes

import (
	"get-shit-done/service"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(authService *service.AuthService) *chi.Mux {
	r := chi.NewRouter()

	// TODO: route for signup

	// TODO: route for login

	return r
}
