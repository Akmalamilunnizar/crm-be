package recurring_invoice

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
)

type AdminRecurringInvoiceHandlerStruct struct {
	service AdminRecurringInvoiceServiceInterface
}

func NewAdminRecurringInvoiceHandler(service AdminRecurringInvoiceServiceInterface) *AdminRecurringInvoiceHandlerStruct {
	return &AdminRecurringInvoiceHandlerStruct{service}
}

func (h AdminRecurringInvoiceHandlerStruct) GetAllRecurringInvoicesHandler(c *fiber.Ctx) error {
	recurringInvoices, err := h.service.GetAllRecurringInvoices()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", recurringInvoices)
}

func (h AdminRecurringInvoiceHandlerStruct) GetRecurringInvoiceByIDHandler(c *fiber.Ctx) error {
	request := IdRecurringInvoiceRequest{}
	request.Id = c.Params("id")

	recurringInvoice, err := h.service.GetRecurringInvoiceByID(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", recurringInvoice)
}

func (h AdminRecurringInvoiceHandlerStruct) CreateRecurringInvoiceHandler(c *fiber.Ctx) error {
	request := CreateRecurringInvoiceRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id").(string)

	recurringInvoice, err := h.service.CreateRecurringInvoice(request, userID)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Recurring Invoice", recurringInvoice)
}

func (h AdminRecurringInvoiceHandlerStruct) UpdateRecurringInvoiceHandler(c *fiber.Ctx) error {
	request := &UpdateRecurringInvoiceRequest{}

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

	recurringInvoice, err := h.service.UpdateRecurringInvoice(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Recurring Invoice", recurringInvoice)
}

func (h AdminRecurringInvoiceHandlerStruct) UpdateRecurringInvoiceStatusHandler(c *fiber.Ctx) error {
	request := &UpdateRecurringInvoiceStatusRequest{}

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

	recurringInvoice, err := h.service.UpdateRecurringInvoiceStatus(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Recurring Invoice Status", recurringInvoice)
}

func (h AdminRecurringInvoiceHandlerStruct) DeleteRecurringInvoiceHandler(c *fiber.Ctx) error {
	request := IdRecurringInvoiceRequest{}
	request.Id = c.Params("id")

	err := h.service.DeleteRecurringInvoice(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Recurring Invoice", nil)
}

func (h AdminRecurringInvoiceHandlerStruct) GenerateInvoiceFromRecurringHandler(c *fiber.Ctx) error {
	request := GenerateInvoiceRequest{}
	request.Id = c.Params("id")

	err := c.BodyParser(&request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	invoice, err := h.service.GenerateInvoiceFromRecurring(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Generate Invoice from Recurring", invoice)
}

func (h AdminRecurringInvoiceHandlerStruct) GetRecurringInvoiceHistoryHandler(c *fiber.Ctx) error {
	request := IdRecurringInvoiceRequest{}
	request.Id = c.Params("id")

	history, err := h.service.GetRecurringInvoiceHistory(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", history)
}

// Manual trigger to process all due recurring invoices
func (h AdminRecurringInvoiceHandlerStruct) RunDueHandler(c *fiber.Ctx) error {
	count, err := h.service.ProcessDueRecurringInvoices()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Processed due recurring invoices", fiber.Map{"generated": count})
}
