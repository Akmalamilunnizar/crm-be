package upload_file

import (
	"fmt"
	"os"
	"path/filepath"
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CommonUploadFileHandlerStruct struct {
	service CommonUploadFileServiceInterface
}

func NewCommonUploadFileHandler(service CommonUploadFileServiceInterface) *CommonUploadFileHandlerStruct {
	return &CommonUploadFileHandlerStruct{service}
}

func (h CommonUploadFileHandlerStruct) GetAllCommonUploadFileHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllCommonUploadFileService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h CommonUploadFileHandlerStruct) GetByIdCommonUploadFileHandler(c *fiber.Ctx) error {
	request := IdCommonUploadFileRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdCommonUploadFileService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h CommonUploadFileHandlerStruct) CreateCommonUploadFileHandler(c *fiber.Ctx) error {
	request := CreateCommonUploadFileRequest{}
	request.File, _ = c.FormFile("file") // Get file separately

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}

	if err := validation.ValidationRequest(&request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(err, ", "), nil)
	}

	extension := filepath.Ext(request.File.Filename)
	if extension != "" {
		extension = extension[1:]
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File has no extension",
		})
	}

	// Normalize legacy path to requested cs-images folder
	if request.Path == "tickets/cs" || strings.HasPrefix(request.Path, "tickets/cs/") {
		request.Path = strings.Replace(request.Path, "tickets/cs", "cs-images", 1)
	}

	// Use configurable base dir; default to ./uploads
	baseDir := os.Getenv("UPLOAD_BASE_DIR")
	if baseDir == "" {
		baseDir = "./uploads"
	}
	// Ensure directory exists
	if mkErr := os.MkdirAll(fmt.Sprintf("%s/%s", baseDir, request.Path), 0755); mkErr != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, mkErr.Error(), nil)
	}
	// Save file
	destination := fmt.Sprintf("%s/%s/%s.%s", baseDir, request.Path, request.Name, extension)
	if err := c.SaveFile(request.File, destination); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	// For trouble ticket CS images, don't create Image records
	// Just return the filename for direct use in trouble_tickets table
	if strings.HasPrefix(request.Path, "cs-images") || strings.HasPrefix(request.Path, "noc-images") {
		// For trouble tickets, just return the filename - don't create Image records
		filename := request.Name + "." + extension
		return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create File", map[string]string{
			"filename": filename,
			"path":     request.Path,
		})
	}

	// For other uploads (customer installations), create Image records as usual
	image, err := h.service.CreateCommonUploadFileService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Area", image)
}

func (h CommonUploadFileHandlerStruct) UpdateCommonUploadFileHandler(c *fiber.Ctx) error {

	request := &UpdateCommonUploadFileRequest{}

	id := c.Params("id")
	err := c.BodyParser(request)
	request.Id = id
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.UpdateCommonUploadFileService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Area", area)
}

func (h CommonUploadFileHandlerStruct) DeleteCommonUploadFileHandler(c *fiber.Ctx) error {
	request := &IdCommonUploadFileRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	image, err := h.service.DeleteCommonUploadFileService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	if err := os.Remove("./public/" + image.FullPath); err != nil {
		if !os.IsNotExist(err) {
			return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
		}
		// If file doesn't exist, proceed to delete the database record
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Area", image)

}
