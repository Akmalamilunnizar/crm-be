package recurring_invoice

import (
	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
)

func AdminRecurringInvoiceRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminRecurringInvoiceRepository(db)
	service := NewAdminRecurringInvoiceService(repository)
	handler := NewAdminRecurringInvoiceHandler(service)

	// Public routes (no auth required)
	app.Get("/:id", handler.GetRecurringInvoiceByIDHandler)
	app.Get("/:id/history", handler.GetRecurringInvoiceHistoryHandler)

	// Protected routes (auth required)
	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllRecurringInvoicesHandler)
	app.Post("/", handler.CreateRecurringInvoiceHandler)
	app.Put("/:id", handler.UpdateRecurringInvoiceHandler)
	app.Put("/:id/status", handler.UpdateRecurringInvoiceStatusHandler)
	app.Post("/:id/generate", handler.GenerateInvoiceFromRecurringHandler)
	app.Delete("/:id", handler.DeleteRecurringInvoiceHandler)

	// Manual batch trigger
	app.Post("/run-due", handler.RunDueHandler)
}
