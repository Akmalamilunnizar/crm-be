package invoice

import (
	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
)

func AdminInvoiceRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminInvoiceRepository(db)
	service := NewAdminInvoiceService(repository)
	handler := NewAdminInvoiceHandler(service)

	app.Get("/:id", handler.GetByIdAdminInvoiceHandler)
	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminInvoiceHandler)
	app.Post("/", handler.CreateAdminInvoiceHandler)
	app.Put("/:id", handler.UpdateAdminInvoiceHandler)
	app.Put("/:id/status", handler.UpdateStatusAdminInvoiceHandler)
	app.Delete("/:id", handler.DeleteAdminInvoiceHandler)
}
