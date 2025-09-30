package network_monitoring

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func NetworkMonitoringRoutes(app fiber.Router) {
	// Apply auth middleware
	app.Use(helpers.VerifyToken)
	
	// Get handler from context
	app.Use(func(c *fiber.Ctx) error {
		handler := c.Locals("networkMonitoringHandler").(*NetworkMonitoringHandler)
		c.Locals("handler", handler)
		return c.Next()
	})

	// Customer device monitoring
	app.Get("/customer/:customerId/status", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.GetCustomerConnectionStatus(c)
	})
	app.Post("/customer/:customerId/monitor", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.MonitorCustomerDevice(c)
	})

	// Netwatch script generation
	app.Post("/netwatch/script", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.GenerateNetwatchScript(c)
	})

	// Troubleshooting
	app.Get("/troubleshooting", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.GetTroubleshootingRecommendations(c)
	})

	// Bulk operations
	app.Post("/bulk/monitor", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.BulkMonitorDevices(c)
	})

	// Statistics
	app.Get("/stats", func(c *fiber.Ctx) error {
		handler := c.Locals("handler").(*NetworkMonitoringHandler)
		return handler.GetMonitoringStats(c)
	})
}
