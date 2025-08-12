package usermanagement

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminUserManagementHandlerStruct struct {
	service AdminUserManagementServiceInterface
}

func NewAdminUserManagementHandler(service *AdminUserManagementServiceStruct) *AdminUserManagementHandlerStruct {
	return &AdminUserManagementHandlerStruct{service}
}

func (h *AdminUserManagementHandlerStruct) GetAllAdminUserManagementHandler(c *fiber.Ctx) error {
	// Implement the logic to get user management details
	request := SearchAdminUserManagementRequest{}
	request.Role = c.Query("role")
	User, err := h.service.GetAllAdminUserManagementService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "User Not Found", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get All User", User)
}

func (h *AdminUserManagementHandlerStruct) GetByIDAdminUserManagementHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	request := IdAdminUserManagementRequest{}
	request.Id = id

	user, err := h.service.GetByIdAdminUserManagementService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "User Not Found", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get User", user)
}

func (h *AdminUserManagementHandlerStruct) CreateAdminUserManagementHandler(c *fiber.Ctx) error {
	request := CreateAdminUserManagementRequest{}

	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)

	}

	user, err := h.service.CreateAdminUserManagementService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get User", user)
}

func (h *AdminUserManagementHandlerStruct) UpdateAdminUserManagementHandler(c *fiber.Ctx) error {
	request := UpdateAdminUserManagementRequest{}

	err := c.BodyParser(&request)
	request.Id = c.Params("id")
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)

	}

	user, err := h.service.UpdateAdminUserManagementService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get User", user)
}

func (h *AdminUserManagementHandlerStruct) DeleteAdminUserManagementHandler(c *fiber.Ctx) error {
	request := IdAdminUserManagementRequest{}

	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)

	}

	user, err := h.service.DeleteAdminUserManagementService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get User", user)
}
