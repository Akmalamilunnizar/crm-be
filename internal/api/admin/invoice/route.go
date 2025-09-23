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
	app.Post("/:id/partial-payment", handler.ProcessPartialPaymentHandler)
	app.Post("/:id/mark-pdf-viewed", handler.MarkPdfViewedHandler)
	app.Post("/print-all-unpaid", handler.PrintAllUnpaidInvoicesHandler)
	app.Delete("/:id", handler.DeleteAdminInvoiceHandler)
}
