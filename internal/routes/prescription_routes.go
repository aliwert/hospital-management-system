package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupPrescriptionRoutes(router fiber.Router) {
	prescriptionRepo := repositories.NewPrescriptionRepository(database.GetSqlDB())
	prescriptionService := services.NewPrescriptionService(prescriptionRepo)
	prescriptionHandler := handlers.NewPrescriptionHandler(prescriptionService)

	prescriptions := router.Group("/prescriptions", middleware.AuthMiddleware())
	{
		prescriptions.Get("/", middleware.RoleMiddleware("admin", "doctor"), prescriptionHandler.GetPrescriptions)
		prescriptions.Get("/:id", prescriptionHandler.GetPrescription)
		prescriptions.Get("/patient/:patient_id", prescriptionHandler.GetPatientPrescriptions)
		prescriptions.Post("/create", middleware.RoleMiddleware("doctor"), prescriptionHandler.CreatePrescription)
		prescriptions.Put("/update/:id", middleware.RoleMiddleware("doctor"), prescriptionHandler.UpdatePrescription)
		prescriptions.Delete("/delete/:id", middleware.RoleMiddleware("admin", "doctor"), prescriptionHandler.DeletePrescription)
	}
}
