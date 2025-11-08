package customerinstallation

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminCustomerInstallationHandlerStruct struct {
	service AdminCustomerInstallationServiceInterface
}

func NewAdminCustomerInstallationHandler(service AdminCustomerInstallationServiceInterface) AdminCustomerInstallationHandlerStruct {
	return AdminCustomerInstallationHandlerStruct{service}
}

func (h AdminCustomerInstallationHandlerStruct) GetAllAdminCustomerInstallationHandler(c *fiber.Ctx) error {
	// Get is_terminal query parameter if provided
	isTerminal := c.Query("is_terminal", "")

	customer, err := h.service.GetAllAdminCustomerInstallationService(isTerminal)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Customer!", customer)
}
func (h AdminCustomerInstallationHandlerStruct) GetByIdAdminCustomerInstallationHandler(c *fiber.Ctx) error {
	request := IdAdminCustomerInstallationRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}
	customer, err := h.service.GetByIdAdminCustomerInstallationService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Id Customer!", customer)
}
func (h AdminCustomerInstallationHandlerStruct) CreateAdminCustomerInstallationHandler(c *fiber.Ctx) error {
	request := CreateAdminCustomerInstallationRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.CreateAdminCustomerInstallationService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Customer!", customer)
}
func (h AdminCustomerInstallationHandlerStruct) UpdateAdminCustomerInstallationHandler(c *fiber.Ctx) error {
	request := UpdateAdminCustomerInstallationRequest{}
	err := c.BodyParser(&request)
	request.Id = c.Params("id")
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.UpdateAdminCustomerInstallationService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Customer!", customer)
}
func (h AdminCustomerInstallationHandlerStruct) DeleteAdminCustomerInstallationHandler(c *fiber.Ctx) error {
	request := IdAdminCustomerInstallationRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	customer, err := h.service.DeleteAdminCustomerInstallationService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Customer!", customer)
}

// NEW: Handler to get installation reports by customer ID
func (h AdminCustomerInstallationHandlerStruct) GetInstallationReportsByCustomerHandler(c *fiber.Ctx) error {
	customerId := c.Params("customerId")

	if customerId == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Customer ID is required", "")
	}

	reports, err := h.service.GetInstallationReportsByCustomerService(customerId)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Installation Reports!", reports)
}
