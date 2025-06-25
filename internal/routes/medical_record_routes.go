package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupMedicalRecordRoutes(router fiber.Router) {
	medicalRecordRepo := repositories.NewMedicalRecordRepository(database.GetSqlDB())
	medicalRecordService := services.NewMedicalRecordService(medicalRecordRepo)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(medicalRecordService)

	records := router.Group("/medical-records", middleware.AuthMiddleware())
	{
		records.Get("/", middleware.RoleMiddleware("admin", "doctor"), medicalRecordHandler.GetMedicalRecords)
		records.Get("/:id", middleware.RoleMiddleware("admin", "doctor"), medicalRecordHandler.GetMedicalRecord)
		records.Get("/patient/:patient_id", middleware.RoleMiddleware("admin", "doctor"), medicalRecordHandler.GetPatientMedicalRecords)
		records.Post("/create", middleware.RoleMiddleware("doctor"), medicalRecordHandler.CreateMedicalRecord)
		records.Put("/update/:id", middleware.RoleMiddleware("doctor"), medicalRecordHandler.UpdateMedicalRecord)
		records.Delete("/delete/:id", middleware.RoleMiddleware("admin"), medicalRecordHandler.DeleteMedicalRecord)
	}
}
