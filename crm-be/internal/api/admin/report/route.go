package report

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminReportRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminReportRepository(db)
	service := NewAdminReportService(repository)
	handler := NewAdminReportHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/internet", handler.GetAllAdminReportHandler)
}
