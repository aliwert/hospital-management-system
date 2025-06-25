package routes

import (
	"github.com/aliwert/go-hospital-management/internal/database"
	"github.com/aliwert/go-hospital-management/internal/handlers"
	"github.com/aliwert/go-hospital-management/internal/middleware"
	"github.com/aliwert/go-hospital-management/internal/repositories"
	"github.com/aliwert/go-hospital-management/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupTestResultRoutes(router fiber.Router) {
	testResultRepo := repositories.NewTestResultRepository(database.GetSqlDB())
	testResultService := services.NewTestResultService(testResultRepo)
	testResultHandler := handlers.NewTestResultHandler(testResultService)

	testResults := router.Group("/test-results", middleware.AuthMiddleware())
	{
		testResults.Post("/", middleware.RoleMiddleware("doctor"), testResultHandler.CreateTestResult)
		testResults.Get("/:id", testResultHandler.GetTestResult)
		testResults.Get("/medical-record/:medical_record_id", testResultHandler.GetTestResultsByMedicalRecord)
		testResults.Put("/update/:id", middleware.RoleMiddleware("doctor"), testResultHandler.UpdateTestResult)
		testResults.Delete("/delete/:id", middleware.RoleMiddleware("doctor"), testResultHandler.DeleteTestResult)
	}
}
