package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupDoctorRoutes(router fiber.Router) {
	// Initialize doctor components
	doctorRepo := repositories.NewDoctorRepository(database.GetSqlDB())
	doctorService := services.NewDoctorService(doctorRepo)
	doctorHandler := handlers.NewDoctorHandler(doctorService)

	// All doctor routes require admin authentication
	doctors := router.Group("/doctors", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		// Admin only routes
		doctors.Get("/", doctorHandler.GetDoctors)
		doctors.Get("/:id", doctorHandler.GetDoctor)
		doctors.Post("/create", doctorHandler.CreateDoctor)
		doctors.Put("/update/:id", doctorHandler.UpdateDoctor)
		doctors.Delete("/delete/:id", doctorHandler.DeleteDoctor)
	}
}
