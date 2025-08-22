package dashboard

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CustomerDashboardHandlerStruct struct {
	service CustomerDashboardServiceInterface
}

func NewCustomerDashboardHandler(service CustomerDashboardServiceInterface) CustomerDashboardHandlerStruct {
	return CustomerDashboardHandlerStruct{service}
}

func (h CustomerDashboardHandlerStruct) MyUserCustomerDashboard(c *fiber.Ctx) error {
	request := c.Locals("user_id").(string)
	myAccount, err := h.service.MyUserCustomerDashboard(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)

	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "My Account", myAccount)
}

func (h CustomerDashboardHandlerStruct) CreatePaymentCustomerDashboard(c *fiber.Ctx) error {

	request := new(CreatePaymentCustomerDashboardRequest)
	if err := c.BodyParser(request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	// Validate the request
	errorMessages := validation.ValidationRequest(request)
	if errorMessages != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errorMessages, ", "), nil)
	}

	SearchInvoice := SearchInvoice{}

	SearchInvoice.UserId = c.Locals("user_id").(string)
	SearchInvoice.InvoiceId = request.InvoiceId

	payment, err := h.service.CreatePaymentCustomerDashboard(SearchInvoice)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "My Account", payment)
}
