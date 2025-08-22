package transaction

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminTrasactionRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminTransactionRepository(db)
	service := NewAdminTransactionService(repository)
	handler := NewAdminTransactionHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminTransactionHandlerStruct)
	app.Get("/:id", handler.GetByIdAdminTransactionHandlerStruct)
	app.Get("/:id/sync", handler.SyncAdminTransactionHandlerStruct)
	app.Post("/", handler.CreateAdminTransactionHandlerStruct)
	app.Put("/:id", handler.UpdateAdminTransactionHandlerStruct)
	app.Delete("/:id", handler.DeleteAdminTransactionHandlerStruct)
}
