package customerinstallation

import (
	"skripsi-be/internal/api/admin/recurring_invoice"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminCustomerInstallationRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminCustomerInstallationRepository(db)
	service := NewAdminCustomerInstallationService(repository)
	handler := NewAdminCustomerInstallationHandler(service)

	// Installation Report routes
	reportRepository := NewAdminInstallationReportRepository(db)
	recurringInvoiceRepo := recurring_invoice.NewAdminRecurringInvoiceRepository(db)
	reportService := NewAdminInstallationReportService(reportRepository, recurringInvoiceRepo)
	reportHandler := NewAdminInstallationReportController(reportService)

	app.Use(helpers.VerifyToken)

	// Installation Report endpoints (must be before /:id route to avoid conflicts)
	app.Get("/report/complete/:id", reportHandler.GetCompleteInstallationReportWithTechnicianPhotos)
	app.Get("/report/complete-view/:id", reportHandler.GetCompleteInstallationReportByView)
	app.Get("/report/technician-team/:id", reportHandler.GetInstallationTechnicianTeam)
	app.Get("/report-complete", reportHandler.GetAllCompleteInstallationReports)
	app.Get("/report/summary/customer", reportHandler.GetInstallationSummaryPerCustomer)
	app.Get("/report/asset/:id", reportHandler.GetInstallationAssetReport)
	app.Get("/report/technician", reportHandler.GetInstallationTechnicianReport)
	app.Post("/report/complete", reportHandler.CreateCompleteInstallationReport)
	app.Put("/report/complete/:id", reportHandler.UpdateCompleteInstallationReport)
	app.Delete("/report/delete/:id", reportHandler.DeleteInstallationReport)

	// Basic CRUD operations
	app.Get("", handler.GetAllAdminCustomerInstallationHandler)
	app.Get("/customer/:customerId", handler.GetInstallationReportsByCustomerHandler) // NEW: Get reports by customer ID
	app.Get("/:id", handler.GetByIdAdminCustomerInstallationHandler)
	app.Post("", handler.CreateAdminCustomerInstallationHandler)
	app.Put("/:id", handler.UpdateAdminCustomerInstallationHandler)
	app.Delete("/:id", handler.DeleteAdminCustomerInstallationHandler)

	// New Installation Report endpoint with multipart form data
	newReportRepository := NewReportInstallationRepository(db)
	newReportService := NewReportInstallationService(newReportRepository, recurringInvoiceRepo)
	newReportHandler := NewReportInstallationController(newReportService)
	app.Post("/report-installations", newReportHandler.CreateReportInstallation)
}
