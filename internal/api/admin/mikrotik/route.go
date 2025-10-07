package mikrotik

import (
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

func MikroTikRoutes(app fiber.Router) {
	// Create handler instance with proper MikroTik service
	handler := NewMikroTikHandler(services.GetSharedMikroTikService())

	// Apply auth middleware
	app.Use(helpers.VerifyToken)

	// Connection management
	app.Post("/connect", handler.Connect)
	app.Post("/disconnect", handler.Disconnect)
	app.Get("/status", handler.GetStatus)

	// Logs and monitoring
	app.Get("/logs", handler.GetLogs)
	app.Get("/logs/realtime", handler.GetRealTimeLogs)

	// Netwatch devices
	app.Get("/netwatch/devices", handler.GetNetwatchDevices)

	// System information
	app.Get("/system/info", handler.GetSystemInfo)

	// Command execution
	app.Post("/execute", handler.ExecuteCommand)

	// Hotspot IP binding management
	app.Post("/hotspot/ip-binding/set-type", handler.SetHotspotIPBindingType)
	app.Get("/hotspot/ip-bindings", handler.GetHotspotIPBindings)
	app.Get("/hotspot/ip-binding/:mac", handler.GetHotspotIPBindingByMAC)

	// DHCP lease management
	app.Post("/dhcp-lease", handler.GetDHCPLease)
}
