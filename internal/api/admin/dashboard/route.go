package dashboard

import (
	"github.com/gofiber/fiber/v2"

	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
)

func AdminDashboardRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminDashboradRepository(db)
	service := NewAdminDashboardService(repository)
	handler := NewAdminDashboardHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("total-income", handler.GetTotalIncome)
	app.Get("total-expenses", handler.GetTotalExpenses)
	app.Get("total-net-worth", handler.GetNetWorth)
	app.Get("total-sales", handler.GetSales)

	app.Get("card-customer", handler.CardCustomer)
	app.Get("card-packet-popular", handler.CardPacketPopular)
	app.Get("card-area-popular", handler.CardAreaPopular)
	app.Get("card-report-cash", handler.CardReportCash)

	// New dashboard endpoints
	app.Get("stats", handler.GetDashboardStats)
	app.Get("recent-invoices", handler.GetRecentInvoices)
	app.Get("recent-transactions", handler.GetRecentTransactions)
	app.Get("customer-growth", handler.GetCustomerGrowth)
	app.Get("revenue-chart", handler.GetRevenueChart)
	app.Get("expenses-chart", handler.GetExpensesChart)
	app.Get("unpaid-customers-chart", handler.GetUnpaidCustomersChart)

	// analityc := app.Group("analytic-graph")
	// analityc.Get("customer")
	// analityc.Get("packet-popular")
	// analityc.Get("area-popular")
	// analityc.Get("report-cash")

}
