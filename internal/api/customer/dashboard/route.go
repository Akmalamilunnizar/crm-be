package dashboard

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func CustomerDashboardRoute(app fiber.Router) {

	db := database.GetDB()
	repository := NewCustomerDashboardRepository(db)
	service := NewCustomerDashboardService(repository)
	handler := NewCustomerDashboardHandler(service)

	// Route
	app.Get("/", helpers.CustomerVerifyToken, handler.MyUserCustomerDashboard)
	app.Post("/", helpers.CustomerVerifyToken, handler.CreatePaymentCustomerDashboard)
}
