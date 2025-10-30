package customer

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminCustomerRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminCustomerRepository(db)
	service := NewAdminCustomerService(repository)
	handler := NewAdminCustomerHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("", handler.GetAllAdminCustomerHandler)
	app.Get("/filter", handler.GetAllAdminCustomerWithFilterHandler)
	app.Get("/:id", handler.GetByIdAdminCustomerHandler)
	app.Get("/:id/detail", handler.GetByIdDetailAdminCustomerHandler)
	app.Get("/:id/related-records", handler.GetCustomerRelatedRecords)
	app.Post("", handler.CreateAdminCustomerHandler)
	app.Put("/:id", handler.UpdateAdminCustomerHandler)
	app.Delete("/:id", handler.DeleteAdminCustomerHandler)
	app.Delete("/:id/with-related", handler.DeleteAdminCustomerWithRelated)

}
