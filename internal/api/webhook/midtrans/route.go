package midtrans

import (
	"skripsi-be/internal/config/database"

	"github.com/gofiber/fiber/v2"
)

func WebhookMidtransRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewWebhookMidtransRepository(db)
	service := NewWebhookMidtransService(repository)
	handler := NewWebhookMidtransHandler(service)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Webhook Midtrans")
	})
	app.Post("/", handler.ReceivedWebhookMidtransHandler)
}
