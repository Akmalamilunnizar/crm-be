package product

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminProductHandlerStruct struct {
	service AdminProductServiceStruct
}

func NewAdminProductHandlerStruct(service AdminProductServiceStruct) AdminProductHandlerStruct {
	return AdminProductHandlerStruct{service}
}

func (h AdminProductHandlerStruct) GetAllAdminProductHandler(c *fiber.Ctx) error {
	products, err := h.service.GetAllAdminProductService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", products)
}

func (h AdminProductHandlerStruct) GetByIdAdminProductHandler(c *fiber.Ctx) error {
	request := IdAdminProductRequest{}
	request.Id = c.Params("id")

	products, err := h.service.GetByIdAdminProductService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", products)
}

func (h AdminProductHandlerStruct) CreateAdminProductHandler(c *fiber.Ctx) error {
	request := CreateAdminProductRequest{}
	err := c.BodyParser(&request)
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	data, err := h.service.CreateAdminProductService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", data)
}

func (h AdminProductHandlerStruct) UpdateAdminProductHandler(c *fiber.Ctx) error {
	request := UpdateAdminProductRequest{}
	request.Id = c.Params("id")

	if request.Id == ":id" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Id Cannot be null", nil)
	}

	err := c.BodyParser(&request)
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	data, err := h.service.UpdateAdminProductService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", data)
}

func (h AdminProductHandlerStruct) DeleteAdminProductHandler(c *fiber.Ctx) error {
	request := IdAdminProductRequest{}
	request.Id = c.Params("id")

	if request.Id == ":id" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Id Cannot be null", nil)
	}

	product, err := h.service.DeleteAdminProductService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", product)
}
