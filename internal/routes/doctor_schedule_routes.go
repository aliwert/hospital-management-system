package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupDoctorScheduleRoutes(router fiber.Router) {
	scheduleRepo := repositories.NewDoctorScheduleRepository(database.GetSqlDB())
	scheduleService := services.NewDoctorScheduleService(scheduleRepo)
	scheduleHandler := handlers.NewDoctorScheduleHandler(scheduleService)

	schedules := router.Group("/doctor-schedules", middleware.AuthMiddleware())
	{
		schedules.Post("/", middleware.RoleMiddleware("admin"), scheduleHandler.CreateSchedule)
		schedules.Get("/doctor/:doctor_id", scheduleHandler.GetDoctorSchedules)
		schedules.Put("/:id", middleware.RoleMiddleware("admin"), scheduleHandler.UpdateSchedule)
		schedules.Delete("/:id", middleware.RoleMiddleware("admin"), scheduleHandler.DeleteSchedule)
	}
}
