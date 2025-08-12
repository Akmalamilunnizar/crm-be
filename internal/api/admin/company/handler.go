package company

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type AdminCompanyHandlerStruct struct {
	service *AdminCompanyServiceStruct
}

func NewAdminCompanyHandler(service *AdminCompanyServiceStruct) *AdminCompanyHandlerStruct {
	return &AdminCompanyHandlerStruct{service}
}

func (h *AdminCompanyHandlerStruct) GetAllAdminCompanyHandler(c *fiber.Ctx) error {
	companies, err := h.service.GetAllAdminCompanyService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, "Get All Admin Company Error", err)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Get All Admin Company", companies)
}

func (h *AdminCompanyHandlerStruct) GetByIdAdminCompanyHandler(c *fiber.Ctx) error {
	companyId := c.Params("id")
	if companyId == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Company ID is required", nil)
	}
	request := &IdAdminCompanyRequest{}
	request.Id = companyId
	companies, err := h.service.GetByIdAdminCompanyService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, "Get Admin Company Error", err)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Get Admin Company", companies)
}

func (h *AdminCompanyHandlerStruct) CreateAdminCompanyHandler(c *fiber.Ctx) error {
	request := &CreateAdminCompanyRequest{}
	err := c.BodyParser(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Error Parsing Data", err)
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Validation Error", err)
	}

	h.service.CreateAdminCompanyService(*request)
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Create Admin Company", nil)

}

func (h *AdminCompanyHandlerStruct) UpdateAdminCompanyHandler(c *fiber.Ctx) error {
	request := &UpdateAdminCompanyRequest{}

	id := c.Params("id")
	err := c.BodyParser(request)
	request.Id = id
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Error Parsing Data", err)
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Validation Error", err)
	}

	h.service.UpdateAdminCompanyService(*request)

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Data Company", "")
}

func (h *AdminCompanyHandlerStruct) DeleteAdminCompanyHandler(c *fiber.Ctx) error {
	request := &IdAdminCompanyRequest{}

	companyId := c.Params("id")
	if companyId == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Company ID is required", nil)
	}

	request.Id = companyId

	h.service.DeleteAdminCompanyService(*request)

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Data Company", "")
}
