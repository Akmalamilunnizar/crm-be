package upload_file

import (
	"skripsi-be/internal/config/database"

	"github.com/gofiber/fiber/v2"
)

func CommonUploadFileRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewCommonUploadFileRepository(db)
	service := NewCommonUploadFileService(repository)
	handler := NewCommonUploadFileHandler(service)

	app.Get("/", handler.GetAllCommonUploadFileHandler)
	app.Post("/", handler.CreateCommonUploadFileHandler)
	app.Delete("/:id", handler.DeleteCommonUploadFileHandler)
}
