package networkdevice

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type AdminNetworkDeviceHandlerStruct struct {
	service AdminNetworkDeviceServiceInterface
}

func NewAdminNetworkDeviceHandler(service AdminNetworkDeviceServiceInterface) *AdminNetworkDeviceHandlerStruct {
	return &AdminNetworkDeviceHandlerStruct{service: service}
}

func (h AdminNetworkDeviceHandlerStruct) GetAllAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	networkDevices, err := h.service.GetAllAdminNetworkDeviceService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get All Network Devices!", networkDevices)
}

func (h AdminNetworkDeviceHandlerStruct) GetByIdAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	request := IdNetworkDeviceRequest{}
	request.Id = c.Params("id")

	networkDevice, err := h.service.GetByIdAdminNetworkDeviceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusNotFound, false, "Network Device not found", "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Network Device!", networkDevice)
}

func (h AdminNetworkDeviceHandlerStruct) GetByCustomerIdAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	customerId := c.Params("customerId")

	networkDevices, err := h.service.GetByCustomerIdAdminNetworkDeviceService(customerId)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Network Devices by Customer!", networkDevices)
}

func (h AdminNetworkDeviceHandlerStruct) CreateAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	request := CreateNetworkDeviceRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	// Validate request
	if validationErrors := validation.ValidationRequest(request); validationErrors != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, validationErrors[0], "")
	}

	networkDevice, err := h.service.CreateAdminNetworkDeviceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusCreated, true, "Success Create Network Device!", networkDevice)
}

func (h AdminNetworkDeviceHandlerStruct) UpdateAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	request := UpdateNetworkDeviceRequest{}
	request.ID = c.Params("id")

	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	// Validate request
	if validationErrors := validation.ValidationRequest(request); validationErrors != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, validationErrors[0], "")
	}

	networkDevice, err := h.service.UpdateAdminNetworkDeviceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Network Device!", networkDevice)
}

func (h AdminNetworkDeviceHandlerStruct) DeleteAdminNetworkDeviceHandler(c *fiber.Ctx) error {
	request := IdNetworkDeviceRequest{}
	request.Id = c.Params("id")

	networkDevice, err := h.service.DeleteAdminNetworkDeviceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Network Device!", networkDevice)
}
