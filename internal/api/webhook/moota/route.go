package moota

import (
	"skripsi-be/internal/config/database"

	"github.com/gofiber/fiber/v2"
)

func WebhookMootaRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewWebhookMootaRepository(db)
	service := NewWebhookMootaService(repository)
	handler := NewWebhookMootaHandler(service)

	app.Get("/", handler.GetAllWebhookMootaHandler)
	app.Get("/:id", handler.GetByIdWebhookMootaHandler)
	app.Post("/", handler.CreateWebhookMootaHandler)
	app.Put("/:id", handler.UpdateWebhookMootaHandler)
	app.Delete("/:id", handler.DeleteWebhookMootaHandler)
}
