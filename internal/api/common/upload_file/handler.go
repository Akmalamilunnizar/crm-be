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

	destination := fmt.Sprintf("./public/uploads/%s/%s.%s", request.Path, request.Name, extension)
	err := c.SaveFile(request.File, destination)
	if err != nil {
		err := os.MkdirAll(fmt.Sprintf("./public/uploads/%s", request.Path), 0755)
		err = c.SaveFile(request.File, destination)
		if err != nil {
			return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
		}
	}

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
