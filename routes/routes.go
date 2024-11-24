package routes

import (
	"get-shit-done/controller"
	"get-shit-done/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(authController *controller.AuthController) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authController.SignUp)
		r.Post("/signin", authController.SignIn)
		r.Post("/refresh", authController.RefreshRefreshToken)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.ValidateAccessToken(authController.JwtService))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("home"))
		})
	})

	return r
}
