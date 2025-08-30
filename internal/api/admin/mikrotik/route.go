package mikrotik

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func MikroTikRoutes(app fiber.Router) {
	// Create handler instance
	handler := NewMikroTikHandler(nil)

	// Apply auth middleware
	app.Use(helpers.VerifyToken)

	// Connection management
	app.Post("/connect", handler.Connect)
	app.Post("/disconnect", handler.Disconnect)
	app.Get("/status", handler.GetStatus)

	// Logs and monitoring
	app.Get("/logs", handler.GetLogs)
	app.Get("/logs/realtime", handler.GetRealTimeLogs)

	// System information
	app.Get("/system/info", handler.GetSystemInfo)

	// Command execution
	app.Post("/execute", handler.ExecuteCommand)
}
