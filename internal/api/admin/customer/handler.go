package customer

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminCustomerHandlerStruct struct {
	service AdminCustomerServiceInterface
}

func NewAdminCustomerHandler(service AdminCustomerServiceInterface) AdminCustomerHandlerStruct {
	return AdminCustomerHandlerStruct{service}
}

func (h AdminCustomerHandlerStruct) GetAllAdminCustomerHandler(c *fiber.Ctx) error {
	customer, err := h.service.GetAllAdminCustomerService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Customer!", customer)
}
func (h AdminCustomerHandlerStruct) GetByIdAdminCustomerHandler(c *fiber.Ctx) error {
	request := IdAdminCustomerRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}
	customer, err := h.service.GetByIdAdminCustomerService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Id Customer!", customer)
}
func (h AdminCustomerHandlerStruct) CreateAdminCustomerHandler(c *fiber.Ctx) error {
	request := CreateAdminCustomerRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.CreateAdminCustomerService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Customer!", customer)
}
func (h AdminCustomerHandlerStruct) UpdateAdminCustomerHandler(c *fiber.Ctx) error {
	request := UpdateAdminCustomerRequest{}
	err := c.BodyParser(&request)
	request.Id = c.Params("id")
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.UpdateAdminCustomerService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Customer!", customer)
}
func (h AdminCustomerHandlerStruct) DeleteAdminCustomerHandler(c *fiber.Ctx) error {
	request := IdAdminCustomerRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.DeleteAdminCustomerService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Customer!", customer)
}
