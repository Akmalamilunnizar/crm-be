package area

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminAreaRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminAreaRepository(db)
	service := NewAdminAreaService(repository)
	handler := NewAdminAreaHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminAreaHandler)
	app.Get("/:id", handler.GetByIdAdminAreaHandler)
	app.Post("/", handler.CreateAdminAreaHandler)
	app.Put("/:id", handler.UpdateAdminAreaHandler)
	app.Delete("/:id", handler.DeleteAdminAreaHandler)
}
