package networkdevice

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminNetworkDeviceRoutes(app *fiber.App, handler *AdminNetworkDeviceHandlerStruct) {
	api := app.Group("/api/admin/network-device")

	// Apply role-based middleware
	api.Use(helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"))

	// Routes
	api.Get("/", handler.GetAllAdminNetworkDeviceHandler)
	api.Get("/:id", handler.GetByIdAdminNetworkDeviceHandler)
	api.Get("/customer/:customerId", handler.GetByCustomerIdAdminNetworkDeviceHandler)
	api.Post("/", handler.CreateAdminNetworkDeviceHandler)
	api.Put("/:id", handler.UpdateAdminNetworkDeviceHandler)
	api.Delete("/:id", handler.DeleteAdminNetworkDeviceHandler)
}
