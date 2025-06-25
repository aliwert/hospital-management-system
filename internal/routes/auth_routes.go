package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	authRepo := repositories.NewUserRepository(database.GetSqlDB())
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := router.Group("/auth")
	{
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
		auth.Post("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	auth.Use(middleware.AuthMiddleware())
	auth.Get("/profile", authHandler.GetProfile)

	// Admin only routes
	admin := auth.Group("/users", middleware.RoleMiddleware("admin"))
	{
		admin.Get("/", authHandler.GetUsers)
		admin.Put("/:id", authHandler.UpdateUser)
		admin.Delete("/:id", authHandler.DeleteUser)
	}
}
