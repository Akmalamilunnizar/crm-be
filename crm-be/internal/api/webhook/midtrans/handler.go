package midtrans

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type WebhookMidtransHandlerStruct struct {
	service WebhookMidtransServiceInterface
}

func NewWebhookMidtransHandler(service WebhookMidtransServiceInterface) *WebhookMidtransHandlerStruct {
	return &WebhookMidtransHandlerStruct{service}
}

func (h WebhookMidtransHandlerStruct) GetAllWebhookMidtransHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllWebhookMidtransService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h WebhookMidtransHandlerStruct) ReceivedWebhookMidtransHandler(c *fiber.Ctx) error {
	request := &MidtransNotification{}
	if err := c.BodyParser(request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	err := h.service.ReceivedWebhookMidtransService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", "Webhook Received")
}
