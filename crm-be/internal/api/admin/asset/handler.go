package asset

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
)

type AdminAssetHandlerStruct struct {
	service AdminAssetServiceInterface
}

func NewAdminAssetHandler(service AdminAssetServiceInterface) *AdminAssetHandlerStruct {
	return &AdminAssetHandlerStruct{service}
}

func (h AdminAssetHandlerStruct) GetAllAdminAssetHandler(c *fiber.Ctx) error {
	assets, err := h.service.GetAllAdminAssetService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", assets)
}

func (h AdminAssetHandlerStruct) GetByIdAdminAssetHandler(c *fiber.Ctx) error {
	request := IdAdminAssetRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	asset, err := h.service.GetByIdAdminAssetService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", asset)
}

func (h AdminAssetHandlerStruct) CreateAdminAssetHandler(c *fiber.Ctx) error {
	request := CreateAdminAssetRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	asset, err := h.service.CreateAdminAssetService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Data", asset)
}

func (h AdminAssetHandlerStruct) UpdateAdminAssetHandler(c *fiber.Ctx) error {
	request := UpdateAdminAssetRequest{}
	request.Id = c.Params("id")
	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	asset, err := h.service.UpdateAdminAssetService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Data", asset)
}

func (h AdminAssetHandlerStruct) DeleteAdminAssetHandler(c *fiber.Ctx) error {
	request := IdAdminAssetRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	asset, err := h.service.DeleteAdminAssetService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Data", asset)
}
