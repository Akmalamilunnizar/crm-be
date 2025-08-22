package role

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminRoleHandlerStruct struct {
	service AdminRoleServiceInterface
}

func NewAdminRoleHandler(service AdminRoleServiceInterface) *AdminRoleHandlerStruct {
	return &AdminRoleHandlerStruct{service}
}

func (h AdminRoleHandlerStruct) GetAllAdminRoleHandler(c *fiber.Ctx) error {
	roles, err := h.service.GetAllAdminRoleService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", roles)
}

func (h AdminRoleHandlerStruct) GetByIdAdminRoleHandler(c *fiber.Ctx) error {
	request := IdAdminRoleRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	role, err := h.service.GetByIdAdminRoleService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", role)
}

func (h AdminRoleHandlerStruct) CreateAdminRoleHandler(c *fiber.Ctx) error {
	request := CreateAdminRoleRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	role, err := h.service.CreateAdminRoleService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", role)
}

func (h AdminRoleHandlerStruct) UpdateAdminRoleHandler(c *fiber.Ctx) error {
	request := UpdateAdminRoleRequest{}
	request.Id = c.Params("id")
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	role, err := h.service.UpdateAdminRoleService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", role)
}

func (h AdminRoleHandlerStruct) DeleteAdminRoleHandler(c *fiber.Ctx) error {
	request := IdAdminRoleRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	role, err := h.service.DeleteAdminRoleService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", role)
}
