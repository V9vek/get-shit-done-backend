package routes

import (
	"get-shit-done/controller"
	"get-shit-done/middleware"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(
	authController *controller.AuthController,
	todoController *controller.TodoController,
) *chi.Mux {
	r := chi.NewRouter()

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_BASE_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Cache duration in seconds
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authController.SignUp)
		r.Post("/signin", authController.SignIn)
		r.Post("/refresh", authController.RefreshRefreshToken)
		r.Get("/verify", authController.VerifyAuth)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.ValidateAccessToken(authController.JwtService))
		r.Get("/", todoController.FindTodoByUserId)
		r.Post("/add", todoController.AddTodo)
		r.Patch("/update/{todoId}", todoController.UpdateTodo)
		r.Post("/logout", authController.Logout)
	})

	return r
}
