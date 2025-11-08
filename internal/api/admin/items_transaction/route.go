package items_transaction

import (
	"skripsi-be/internal/config/database"

	"github.com/gofiber/fiber/v2"
)

func ItemsTransactionRoute(router fiber.Router) {
	// Initialize repository, service, and handler
	repository := NewItemsTransactionRepository(database.GetDB())
	service := NewItemsTransactionService(repository)
	handler := NewItemsTransactionHandler(service)

	// Apply auth middleware if needed (assuming you have one)
	// router.Use(authMiddleware)

	// Routes
	router.Get("/", handler.GetAllItemsTransactionHandler)
	router.Get("/:id", handler.GetByIdItemsTransactionHandler)
	router.Post("/", handler.CreateItemsTransactionHandler)
	router.Put("/:id", handler.UpdateItemsTransactionHandler)
	router.Delete("/:id", handler.DeleteItemsTransactionHandler)
}

