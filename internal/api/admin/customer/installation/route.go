package customerinstallation

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminCustomerInstallationRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminCustomerInstallationRepository(db)
	service := NewAdminCustomerInstallationService(repository)
	handler := NewAdminCustomerInstallationHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("", handler.GetAllAdminCustomerInstallationHandler)
	app.Get("/:id", handler.GetByIdAdminCustomerInstallationHandler)
	app.Post("", handler.CreateAdminCustomerInstallationHandler)
	app.Put("/:id", handler.UpdateAdminCustomerInstallationHandler)
	app.Delete("/:id", handler.DeleteAdminCustomerInstallationHandler)

}
