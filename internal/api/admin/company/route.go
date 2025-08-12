package company

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminCompanyRoute(app fiber.Router) {
	dbmysql := database.GetDB()
	repository := NewAdminCompanyRepository(dbmysql)
	service := NewAdminCompanyService(repository)
	handler := NewAdminCompanyHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminCompanyHandler)
	app.Get("/:id", handler.GetByIdAdminCompanyHandler)
	app.Post("/", handler.CreateAdminCompanyHandler)
	app.Put("/:id", handler.UpdateAdminCompanyHandler)
	app.Delete("/:id", handler.DeleteAdminCompanyHandler)
}
