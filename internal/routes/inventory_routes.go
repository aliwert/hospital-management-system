package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupInventoryRoutes(router fiber.Router) {
	// Initialize inventory components
	inventoryRepo := repositories.NewInventoryRepository(database.GetSqlDB())
	inventoryService := services.NewInventoryService(inventoryRepo)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Inventory routes (all require admin authentication)
	inventory := router.Group("/inventory", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		inventory.Get("/", inventoryHandler.GetAllInventory)
		inventory.Get("/:id", inventoryHandler.GetInventory)
		inventory.Post("/create", inventoryHandler.CreateInventory)
		inventory.Put("/:id", inventoryHandler.UpdateInventory)
		inventory.Delete("/:id", inventoryHandler.DeleteInventory)
	}
}
