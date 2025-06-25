package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupAppointmentRoutes(router fiber.Router) {
	appointmentRepo := repositories.NewAppointmentRepository(database.GetSqlDB())
	appointmentService := services.NewAppointmentService(appointmentRepo)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	appointments := router.Group("/appointments", middleware.AuthMiddleware())
	{
		appointments.Post("/", middleware.RoleMiddleware("patient"), appointmentHandler.CreateAppointment)
		appointments.Get("/", middleware.RoleMiddleware("admin", "doctor"), appointmentHandler.GetAppointments)
		appointments.Get("/:id", appointmentHandler.GetAppointment)
		appointments.Put("/:id", middleware.RoleMiddleware("doctor"), appointmentHandler.UpdateAppointment)
		appointments.Delete("/:id", middleware.RoleMiddleware("admin"), appointmentHandler.DeleteAppointment)
	}
}
