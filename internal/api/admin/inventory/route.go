package inventory

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func InventoryRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewInventoryRepository(db)
	service := NewInventoryService(repository)
	handler := NewInventoryHandler(service)

	// Apply authentication middleware
	app.Use(helpers.VerifyToken)

	// Purchase routes
	app.Post("/purchases", handler.CreatePurchaseHandler)

	// Deployment routes
	app.Post("/deployments", handler.CreateDeploymentHandler)

	// Inventory status routes
	app.Get("/inventory/status", handler.GetInventoryStatusHandler)
}

