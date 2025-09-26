package networkdevice

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminNetworkDeviceRoutes(admin fiber.Router, handler *AdminNetworkDeviceHandlerStruct) {
	api := admin.Group("/network-device")

	// Apply authentication middleware first
	api.Use(helpers.VerifyToken)
	// Apply role-based middleware
	api.Use(helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "TECHNICIAN"))

	// Routes
	api.Get("/", handler.GetAllAdminNetworkDeviceHandler)
	api.Get("/:id", handler.GetByIdAdminNetworkDeviceHandler)
	api.Get("/customer/:customerId", handler.GetByCustomerIdAdminNetworkDeviceHandler)
	api.Post("/", handler.CreateAdminNetworkDeviceHandler)
	api.Put("/:id", handler.UpdateAdminNetworkDeviceHandler)
	api.Delete("/:id", handler.DeleteAdminNetworkDeviceHandler)
}
