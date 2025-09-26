package networkdevice

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminNetworkDeviceRoutes(admin fiber.Router, handler *AdminNetworkDeviceHandlerStruct) {
	api := admin.Group("/network-device")

	// Routes with authentication and authorization
	api.Get("/", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.GetAllAdminNetworkDeviceHandler)
	api.Get("/:id", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.GetByIdAdminNetworkDeviceHandler)
	api.Get("/customer/:customerId", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.GetByCustomerIdAdminNetworkDeviceHandler)
	api.Post("/", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.CreateAdminNetworkDeviceHandler)
	api.Put("/:id", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.UpdateAdminNetworkDeviceHandler)
	api.Delete("/:id", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), handler.DeleteAdminNetworkDeviceHandler)
}
