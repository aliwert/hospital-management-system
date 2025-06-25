package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupPatientRoutes(router fiber.Router) {
	patientRepo := repositories.NewPatientRepository(database.GetSqlDB())
	patientService := services.NewPatientService(patientRepo)
	patientHandler := handlers.NewPatientHandler(patientService)

	patients := router.Group("/patients", middleware.AuthMiddleware())
	{
		patients.Post("/", middleware.RoleMiddleware("admin"), patientHandler.CreatePatient)
		patients.Get("/", middleware.RoleMiddleware("admin", "doctor"), patientHandler.GetPatients)
		patients.Get("/:id", middleware.RoleMiddleware("admin", "doctor"), patientHandler.GetPatient)
		patients.Put("/:id", middleware.RoleMiddleware("admin"), patientHandler.UpdatePatient)
		patients.Delete("/:id", middleware.RoleMiddleware("admin"), patientHandler.DeletePatient)
	}
}
