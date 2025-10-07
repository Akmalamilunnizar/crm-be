package asset_transaction

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AssetTransactionRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAssetTransactionRepository(db)
	service := NewAssetTransactionService(repository)
	handler := NewAssetTransactionHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAssetTransactionHandler)
	app.Get("/:id", handler.GetByIdAssetTransactionHandler)
	app.Post("/", handler.CreateAssetTransactionHandler)
	app.Put("/:id", handler.UpdateAssetTransactionHandler)
	app.Delete("/:id", handler.DeleteAssetTransactionHandler)
}

