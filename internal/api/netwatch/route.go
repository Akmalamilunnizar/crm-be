package netwatchapi

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func NetwatchRoutes(app *fiber.App, handler *Handler) {
	// Netwatch API group
	netwatchGroup := app.Group("/api/netwatch")

	// Device management routes
	devices := netwatchGroup.Group("/devices")
	devices.Get("/", handler.ListDevices)
	devices.Get("/:id", handler.GetDevice)
	devices.Post("/", helpers.RequireRoles("ADMIN", "NOC"), handler.CreateDevice)
	devices.Put("/:id", helpers.RequireRoles("ADMIN", "NOC"), handler.UpdateDevice)
	devices.Delete("/:id", helpers.RequireRoles("ADMIN"), handler.DeleteDevice)

	// Event management routes
	events := netwatchGroup.Group("/events")
	events.Get("/", handler.ListEvents)
	events.Get("/:id", handler.GetEvent)
	events.Post("/", handler.CreateEvent) // No auth required for webhook/syslog

	// Configuration management routes
	configs := netwatchGroup.Group("/configs")
	configs.Get("/", helpers.RequireRoles("ADMIN"), handler.ListConfigs)
	configs.Post("/", helpers.RequireRoles("ADMIN"), handler.CreateConfig)
	configs.Put("/:id", helpers.RequireRoles("ADMIN"), handler.UpdateConfig)
	configs.Delete("/:id", helpers.RequireRoles("ADMIN"), handler.DeleteConfig)

	// Sync and monitoring routes
	netwatchGroup.Get("/sync", helpers.RequireRoles("ADMIN", "NOC"), handler.SyncDevices)
	netwatchGroup.Get("/status", handler.GetMonitoringStatus)
}
