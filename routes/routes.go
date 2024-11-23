package routes

import (
	"get-shit-done/controller"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(authController *controller.AuthController) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authController.SignUp)
		r.Post("/signin", authController.SignIn)
		r.Post("/refresh", authController.RefreshRefreshToken)
	})

	return r
}
