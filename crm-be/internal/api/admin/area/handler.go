package area

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminAreaHandlerStruct struct {
	service AdminAreaServiceInterface
}

func NewAdminAreaHandler(service AdminAreaServiceInterface) *AdminAreaHandlerStruct {
	return &AdminAreaHandlerStruct{service}
}

func (h AdminAreaHandlerStruct) GetAllAdminAreaHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllAdminAreaService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminAreaHandlerStruct) GetByIdAdminAreaHandler(c *fiber.Ctx) error {
	request := IdAdminAreaRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdAdminAreaService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminAreaHandlerStruct) CreateAdminAreaHandler(c *fiber.Ctx) error {
	request := CreateAdminAreaRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.CreateAdminAreaService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Area", area)
}

func (h AdminAreaHandlerStruct) UpdateAdminAreaHandler(c *fiber.Ctx) error {

	request := &UpdateAdminAreaRequest{}

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

	area, err := h.service.UpdateAdminAreaService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Area", area)
}

func (h AdminAreaHandlerStruct) DeleteAdminAreaHandler(c *fiber.Ctx) error {
	request := &IdAdminAreaRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.DeleteAdminAreaService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Area", area)

}
