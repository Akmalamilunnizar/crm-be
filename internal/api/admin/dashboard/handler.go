package dashboard

import (
	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/helpers"
)

type AdminDashboardHandlerStruct struct {
	service AdminDashboardServiceInterface
}

func NewAdminDashboardHandler(service AdminDashboardServiceInterface) AdminDashboardHandlerStruct {
	return AdminDashboardHandlerStruct{service}
}

func (h AdminDashboardHandlerStruct) GetTotalIncome(c *fiber.Ctx) error {
	type GetTotalIncome struct {
		TotalIncome int64 `json:"total_income"`
	}

	totalIncome, err := h.service.GetTotalIncome()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Total Income", nil)
	}

	data := GetTotalIncome{
		TotalIncome: totalIncome,
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) GetTotalExpenses(c *fiber.Ctx) error {
	type GetTotalExpenses struct {
		TotalExpenses int64 `json:"total_expenses"`
	}

	totalExpenses, err := h.service.GetTotalExpenses()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Total Income", nil)
	}
	data := GetTotalExpenses{
		TotalExpenses: totalExpenses,
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) GetNetWorth(c *fiber.Ctx) error {
	type GetNetWorth struct {
		TotalNetWorth int64 `json:"total_net_worth"`
	}

	data := GetNetWorth{
		TotalNetWorth: 1000,
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) GetSales(c *fiber.Ctx) error {
	type GetSales struct {
		TotalSale int64 `json:"total_sales"`
	}

	totalCustomer, err := h.service.GetTotalCustomer()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Total Customer", nil)
	}

	data := GetSales{
		TotalSale: totalCustomer,
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) CardCustomer(c *fiber.Ctx) error {
	type CardCustomer struct {
		TotalCutomer int64                  `json:"total_sales"`
		DataGraph    map[string]interface{} `json:"data_graph"`
	}
	data, err := h.service.CardCustomer()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Card", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) CardPacketPopular(c *fiber.Ctx) error {
	type CardPacketPopular struct {
		Total     int64                  `json:"total_sales"`
		DataGraph map[string]interface{} `json:"data_graph"`
	}
	data, err := h.service.CardPacketPopular()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Card", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) CardAreaPopular(c *fiber.Ctx) error {
	data, err := h.service.CardAreaPopular()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Card", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}

func (h AdminDashboardHandlerStruct) CardReportCash(c *fiber.Ctx) error {
	data, err := h.service.CardReportCash()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Card", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Data", data)
}
