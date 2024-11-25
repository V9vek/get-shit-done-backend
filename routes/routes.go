package routes

import (
	"get-shit-done/controller"
	"get-shit-done/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(
	authController *controller.AuthController,
	todoController *controller.TodoController,
) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authController.SignUp)
		r.Post("/signin", authController.SignIn)
		r.Post("/refresh", authController.RefreshRefreshToken)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.ValidateAccessToken(authController.JwtService))
		r.Get("/", todoController.FindTodoByUserId)
		r.Post("/add", todoController.AddTodo)
		r.Patch("/update/{todoId}", todoController.AddTodo)
	})

	return r
}
