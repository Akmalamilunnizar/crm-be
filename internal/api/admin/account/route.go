package account

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminAccountRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminAccountRepository(db)
	service := NewAdminAccountService(repository)
	handler := NewAdminAccountHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminAccountHandler)
	app.Get("/:id", handler.GetByIdAdminAccountHandler)
	app.Post("/", handler.CreateAdminAccountHandler)
	app.Put("/:id", handler.UpdateAdminAccountHandler)
	app.Delete("/:id", handler.DeleteAdminAccountHandler)
}
