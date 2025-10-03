package dashboard

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/helpers"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

	// Calculate net worth: Total Income - Total Expenses
	totalIncome, err := h.service.GetTotalIncome()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Total Income", nil)
	}

	totalExpenses, err := h.service.GetTotalExpenses()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Data Total Expenses", nil)
	}

	netWorth := totalIncome - totalExpenses

	data := GetNetWorth{
		TotalNetWorth: netWorth,
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

func (h AdminDashboardHandlerStruct) GetDashboardStats(c *fiber.Ctx) error {
	data, err := h.service.GetDashboardStats()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Dashboard Stats", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Dashboard Stats", data)
}

func (h AdminDashboardHandlerStruct) GetRecentInvoices(c *fiber.Ctx) error {
	data, err := h.service.GetRecentInvoices()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Recent Invoices", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Recent Invoices", map[string]interface{}{
		"invoices": data,
	})
}

func (h AdminDashboardHandlerStruct) GetRecentTransactions(c *fiber.Ctx) error {
	data, err := h.service.GetRecentTransactions()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Recent Transactions", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Recent Transactions", map[string]interface{}{
		"transactions": data,
	})
}

func (h AdminDashboardHandlerStruct) GetCustomerGrowth(c *fiber.Ctx) error {
	// Optional query params: days, year_start, year_end
	days := c.QueryInt("days", 0)
	yStart := c.QueryInt("year_start", 0)
	yEnd := c.QueryInt("year_end", 0)

	var startPtr *time.Time
	var endPtr *time.Time
	if yStart != 0 && yEnd != 0 {
		// Year range has priority
		start := time.Date(min(yStart, yEnd), 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(max(yStart, yEnd), 12, 31, 0, 0, 0, 0, time.UTC)
		startPtr = &start
		endPtr = &end
	} else if days > 0 {
		end := time.Now().Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -days+1)
		startPtr = &start
		endPtr = &end
	}

	data, err := h.service.GetCustomerGrowth(startPtr, endPtr)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Customer Growth", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Customer Growth", data)
}

func (h AdminDashboardHandlerStruct) GetRevenueChart(c *fiber.Ctx) error {
	days := c.QueryInt("days", 0)
	yStart := c.QueryInt("year_start", 0)
	yEnd := c.QueryInt("year_end", 0)

	var startPtr *time.Time
	var endPtr *time.Time
	if yStart != 0 && yEnd != 0 {
		start := time.Date(min(yStart, yEnd), 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(max(yStart, yEnd), 12, 31, 0, 0, 0, 0, time.UTC)
		startPtr = &start
		endPtr = &end
	} else if days > 0 {
		end := time.Now().Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -days+1)
		startPtr = &start
		endPtr = &end
	}

	data, err := h.service.GetRevenueChart(startPtr, endPtr)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Revenue Chart", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Revenue Chart", data)
}

func (h AdminDashboardHandlerStruct) GetExpensesChart(c *fiber.Ctx) error {
	days := c.QueryInt("days", 0)
	yStart := c.QueryInt("year_start", 0)
	yEnd := c.QueryInt("year_end", 0)

	var startPtr *time.Time
	var endPtr *time.Time
	if yStart != 0 && yEnd != 0 {
		start := time.Date(min(yStart, yEnd), 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(max(yStart, yEnd), 12, 31, 0, 0, 0, 0, time.UTC)
		startPtr = &start
		endPtr = &end
	} else if days > 0 {
		end := time.Now().Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -days+1)
		startPtr = &start
		endPtr = &end
	}

	data, err := h.service.GetExpensesChart(startPtr, endPtr)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Expenses Chart", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Expenses Chart", data)
}

func (h AdminDashboardHandlerStruct) GetUnpaidCustomersChart(c *fiber.Ctx) error {
	days := c.QueryInt("days", 0)
	yStart := c.QueryInt("year_start", 0)
	yEnd := c.QueryInt("year_end", 0)

	var startPtr *time.Time
	var endPtr *time.Time
	if yStart != 0 && yEnd != 0 {
		start := time.Date(min(yStart, yEnd), 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(max(yStart, yEnd), 12, 31, 0, 0, 0, 0, time.UTC)
		startPtr = &start
		endPtr = &end
	} else if days > 0 {
		end := time.Now().Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -days+1)
		startPtr = &start
		endPtr = &end
	}

	data, err := h.service.GetUnpaidCustomersChart(startPtr, endPtr)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Unpaid Customers Chart", nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Unpaid Customers Chart", data)
}

func (h AdminDashboardHandlerStruct) GetUnpaidCustomersList(c *fiber.Ctx) error {
	fmt.Println("🔧 [DEBUG] GetUnpaidCustomersList: Handler called")

	data, err := h.service.GetUnpaidCustomersList()
	if err != nil {
		fmt.Printf("🔧 [DEBUG] GetUnpaidCustomersList: Service error - %v\n", err)
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Failed Get Unpaid Customers List", nil)
	}

	fmt.Printf("🔧 [DEBUG] GetUnpaidCustomersList: Success - found %d unpaid customers\n", len(data))
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success Get Unpaid Customers List", map[string]interface{}{
		"unpaid_customers": data,
	})
}
