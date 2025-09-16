package wa

import (
	"github.com/gofiber/fiber/v2"
)

func WARoutes(app *fiber.App) {
	handler := NewWAHandler()

	// WhatsApp message sending endpoint
	app.Post("/send-message", handler.SendMessage)
}
