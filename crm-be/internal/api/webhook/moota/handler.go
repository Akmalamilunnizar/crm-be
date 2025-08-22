package moota

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type WebhookMootaHandlerStruct struct {
	service WebhookMootaServiceInterface
}

func NewWebhookMootaHandler(service WebhookMootaServiceInterface) *WebhookMootaHandlerStruct {
	return &WebhookMootaHandlerStruct{service}
}

func (h WebhookMootaHandlerStruct) GetAllWebhookMootaHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllWebhookMootaService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h WebhookMootaHandlerStruct) GetByIdWebhookMootaHandler(c *fiber.Ctx) error {
	request := IdWebhookMootaRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdWebhookMootaService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h WebhookMootaHandlerStruct) CreateWebhookMootaHandler(c *fiber.Ctx) error {
	request := []CreateWebhookMootaRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	for _, vrequest := range request {
		errValidation := validation.ValidationRequest(&vrequest)
		if errValidation != nil {
			return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
		}
	}

	area, err := h.service.CreateWebhookMootaService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Area", area)
}

func (h WebhookMootaHandlerStruct) UpdateWebhookMootaHandler(c *fiber.Ctx) error {

	request := &UpdateWebhookMootaRequest{}

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

	area, err := h.service.UpdateWebhookMootaService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Area", area)
}

func (h WebhookMootaHandlerStruct) DeleteWebhookMootaHandler(c *fiber.Ctx) error {
	request := &IdWebhookMootaRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.DeleteWebhookMootaService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Area", area)

}
