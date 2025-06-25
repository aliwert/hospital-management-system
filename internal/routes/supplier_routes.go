package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupSupplierRoutes(router fiber.Router) {
	// Initialize supplier components
	supplierRepo := repositories.NewSupplierRepository(database.GetSqlDB())
	supplierService := services.NewSupplierService(supplierRepo)
	supplierHandler := handlers.NewSupplierHandler(supplierService)

	// Supplier routes (all require admin authentication)
	suppliers := router.Group("/suppliers", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		suppliers.Get("/", supplierHandler.GetSuppliers)
		suppliers.Get("/:id", supplierHandler.GetSupplier)
		suppliers.Post("/create", supplierHandler.CreateSupplier)
		suppliers.Put("/update/:id", supplierHandler.UpdateSupplier)
		suppliers.Delete("/delete/:id", supplierHandler.DeleteSupplier)
	}
}
