package invoice

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
)

type AdminInvoiceHandlerStruct struct {
	service AdminInvoiceServiceInterface
}

func NewAdminInvoiceHandler(service AdminInvoiceServiceInterface) *AdminInvoiceHandlerStruct {
	return &AdminInvoiceHandlerStruct{service}
}

func (h AdminInvoiceHandlerStruct) GetAllAdminInvoiceHandler(c *fiber.Ctx) error {
	areas, err := h.service.GetAllAdminInvoiceService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminInvoiceHandlerStruct) GetByIdAdminInvoiceHandler(c *fiber.Ctx) error {
	request := IdAdminInvoiceRequest{}
	request.Id = c.Params("id")
	areas, err := h.service.GetByIdAdminInvoiceService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", areas)
}

func (h AdminInvoiceHandlerStruct) CreateAdminInvoiceHandler(c *fiber.Ctx) error {
	request := CreateAdminInvoiceRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.CreateAdminInvoiceService(request)

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	// If backend attached warning/error context (e.g., Mikrotik scheduler not found), surface as HTTP 400
	if area.Link != "" && (strings.Contains(area.Link, "MikroTik") || strings.Contains(strings.ToLower(area.Link), "scheduler")) {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, area.Link, area)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Create Invoice", area)
}

func (h AdminInvoiceHandlerStruct) UpdateAdminInvoiceHandler(c *fiber.Ctx) error {

	request := &UpdateAdminInvoiceRequest{}

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

	area, err := h.service.UpdateAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	if area.Link != "" && (strings.Contains(area.Link, "MikroTik") || strings.Contains(strings.ToLower(area.Link), "scheduler")) {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, area.Link, area)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Invoice", area)
}

func (h AdminInvoiceHandlerStruct) UpdateStatusAdminInvoiceHandler(c *fiber.Ctx) error {

	request := &UpdateStatusAdminInvoiceRequest{}

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

	area, err := h.service.UpdateStatusAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	if area.Link != "" && (strings.Contains(area.Link, "MikroTik") || strings.Contains(strings.ToLower(area.Link), "scheduler")) {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, area.Link, area)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Update Invoice", area)
}

func (h AdminInvoiceHandlerStruct) DeleteAdminInvoiceHandler(c *fiber.Ctx) error {
	request := &IdAdminInvoiceRequest{}
	request.Id = c.Params("id")
	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	area, err := h.service.DeleteAdminInvoiceService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Delete Invoice", area)

}

func (h AdminInvoiceHandlerStruct) ProcessPartialPaymentHandler(c *fiber.Ctx) error {
	request := &PartialPaymentRequest{}
	request.Id = c.Params("id")

	err := c.BodyParser(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	invoice, err := h.service.ProcessPartialPaymentService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Partial payment processed successfully", invoice)
}

func (h AdminInvoiceHandlerStruct) MarkPdfViewedHandler(c *fiber.Ctx) error {
	request := &IdAdminInvoiceRequest{}
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	invoice, err := h.service.MarkPdfViewedService(*request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "PDF marked as viewed", invoice)
}

// PrintAllUnpaidInvoicesHandler prints all unpaid invoices to thermal printer
func (h AdminInvoiceHandlerStruct) PrintAllUnpaidInvoicesHandler(c *fiber.Ctx) error {
	result, err := h.service.PrintAllUnpaidInvoicesService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "All unpaid invoices printed successfully", result)
}

// GetRouterJobsByInvoice returns all router jobs for an invoice
func (h AdminInvoiceHandlerStruct) GetRouterJobsByInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "missing id", nil)
	}
	db := database.GetDB()
	var jobs []entities.RouterJob
	if err := db.Where("invoice_id = ?", id).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", jobs)
}

// RetryRouterJobsByInvoice resets failed jobs to pending
func (h AdminInvoiceHandlerStruct) RetryRouterJobsByInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "missing id", nil)
	}
	db := database.GetDB()
	// reset error jobs and schedule soon
	now := time.Now().Add(2 * time.Second)
	if err := db.Model(&entities.RouterJob{}).
		Where("invoice_id = ? AND status = ?", id, entities.RouterJobStatusError).
		Updates(map[string]interface{}{"status": entities.RouterJobStatusPending, "next_run_at": now, "updated_at": time.Now(), "last_error": nil}).Error; err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Jobs retried", nil)
}
