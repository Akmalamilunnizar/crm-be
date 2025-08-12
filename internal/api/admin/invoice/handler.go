package invoice

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
)

type AdminInvoiceHandlerStruct struct {
	service AdminInvoiceServiceInterface
}

func NewAdminInvoiceHandler(service AdminInvoiceServiceInterface) *AdminInvoiceHandlerStruct {
	return &AdminInvoiceHandlerStruct{service}
}

func (h AdminInvoiceHandlerStruct) GetAllAdminInvoiceHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllAdminInvoiceService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminInvoiceHandlerStruct) GetByIdAdminInvoiceHandler(c *fiber.Ctx) error {
	request := IdAdminInvoiceRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdAdminInvoiceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminInvoiceHandlerStruct) CreateAdminInvoiceHandler(c *fiber.Ctx) error {
	request := CreateAdminInvoiceRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.CreateAdminInvoiceService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Invoice", area)
}

func (h AdminInvoiceHandlerStruct) UpdateAdminInvoiceHandler(c *fiber.Ctx) error {

	request := &UpdateAdminInvoiceRequest{}

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

	area, err := h.service.UpdateAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Invoice", area)
}

func (h AdminInvoiceHandlerStruct) UpdateStatusAdminInvoiceHandler(c *fiber.Ctx) error {

	request := &UpdateStatusAdminInvoiceRequest{}

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

	area, err := h.service.UpdateStatusAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Invoice", area)
}

func (h AdminInvoiceHandlerStruct) DeleteAdminInvoiceHandler(c *fiber.Ctx) error {
	request := &IdAdminInvoiceRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.DeleteAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Invoice", area)

}
