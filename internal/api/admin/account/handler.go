package account

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminAccountHandlerStruct struct {
	service AdminAccountServiceInterface
}

func NewAdminAccountHandler(service AdminAccountServiceInterface) *AdminAccountHandlerStruct {
	return &AdminAccountHandlerStruct{service}
}

func (h AdminAccountHandlerStruct) GetAllAdminAccountHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllAdminAccountService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminAccountHandlerStruct) GetByIdAdminAccountHandler(c *fiber.Ctx) error {
	request := IdAdminAccountRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdAdminAccountService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminAccountHandlerStruct) CreateAdminAccountHandler(c *fiber.Ctx) error {
	request := CreateAdminAccountRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.CreateAdminAccountService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Area", area)
}

func (h AdminAccountHandlerStruct) UpdateAdminAccountHandler(c *fiber.Ctx) error {

	request := &UpdateAdminAccountRequest{}

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

	area, err := h.service.UpdateAdminAccountService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Area", area)
}

func (h AdminAccountHandlerStruct) DeleteAdminAccountHandler(c *fiber.Ctx) error {
	request := &IdAdminAccountRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.DeleteAdminAccountService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Area", area)

}
