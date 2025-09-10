package monitoring

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RouteCustomerMonitoring(app *fiber.App) {
	// Initialize dependencies
	db := database.GetDB()
	repo := NewRepo(db)
	service := NewService(db)
	handler := NewHandler(service)

	// Customer monitoring routes (requires admin authentication for summary and customers list)
	customerMonitoring := app.Group("/api/customer/monitoring", middleware.AdminAuthMiddleware)

	// Connection history endpoint (for individual customer)
	customerMonitoring.Get("/history", handler.GetConnectionHistory)

	// WebSocket endpoint for real-time monitoring
	customerMonitoring.Get("/ws", websocket.New(WebSocketHandler))

	// Health check endpoint
	customerMonitoring.Get("/health", func(c *fiber.Ctx) error {
		return helpers.ResponseUtils(c, 200, true, "Customer monitoring service is healthy", map[string]interface{}{
			"service": "customer-monitoring",
			"version": "1.0.0",
			"status":  "active",
		})
	})
}
