package telegramapi

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func TelegramRoutes(app fiber.Router) {
	handler := NewHandler()

	// Telegram test endpoints
	g := app.Group("/telegram")

	// Test basic notification
	g.Post("/test", helpers.VerifyToken, helpers.RequireRoles("ADMIN"), handler.TestNotification)

	// Test ticket notification
	g.Post("/test-ticket", helpers.VerifyToken, helpers.RequireRoles("ADMIN"), handler.TestTicketNotification)

	// Get channel information
	g.Get("/channels", helpers.VerifyToken, helpers.RequireRoles("ADMIN"), handler.GetChannelInfo)
}
