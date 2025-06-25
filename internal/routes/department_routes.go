package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupDepartmentRoutes(router fiber.Router) {
	departmentRepo := repositories.NewDepartmentRepository(database.GetSqlDB())
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	departments := router.Group("/departments", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		departments.Get("/", departmentHandler.GetDepartments)
		departments.Get("/:id", departmentHandler.GetDepartment)
		departments.Post("/", departmentHandler.CreateDepartment)
		departments.Put("/:id", departmentHandler.UpdateDepartment)
		departments.Delete("/:id", departmentHandler.DeleteDepartment)
	}
}
