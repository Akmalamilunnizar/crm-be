package routes

import (
	"skripsi-be/internal/api/admin/account"
	"skripsi-be/internal/api/admin/area"
	"skripsi-be/internal/api/admin/asset"
	"skripsi-be/internal/api/admin/company"
	"skripsi-be/internal/api/admin/customer"
	customerinstallation "skripsi-be/internal/api/admin/customer/installation"
	"skripsi-be/internal/api/admin/dashboard"
	"skripsi-be/internal/api/admin/invoice"
	"skripsi-be/internal/api/admin/product"
	"skripsi-be/internal/api/admin/report"
	"skripsi-be/internal/api/admin/role"
	"skripsi-be/internal/api/admin/transaction"
	usermanagement "skripsi-be/internal/api/admin/user-management"
	authapi "skripsi-be/internal/api/auth"
	upload_file "skripsi-be/internal/api/common/upload_file"
	customerdashboard "skripsi-be/internal/api/customer/dashboard"
	midtrans "skripsi-be/internal/api/webhook/midtrans"
	"skripsi-be/internal/api/webhook/moota"

	"github.com/gofiber/fiber/v2"
)

func RouteFiber(app *fiber.App) {
	app.Static("/", "./public")

	api := app.Group("/api")

	upload_file.CommonUploadFileRoute(api.Group("/file-upload"))

	auth := api.Group("/auth")
	authapi.AuthRoute(auth)

	admin := api.Group("/admin")
	dashboard.AdminDashboardRoute(admin.Group("/dashboard"))
	company.AdminCompanyRoute(admin.Group("/company"))
	account.AdminAccountRoute(admin.Group("/account"))
	usermanagement.AdminUserManagementRoute(admin.Group("/user-management"))
	role.AdminRoleRoute(admin.Group("/role"))
	product.AdminProductRoute(admin.Group("/product"))
	report.AdminReportRoute(admin.Group("/report"))
	area.AdminAreaRoute(admin.Group("/area"))
	customer.AdminCustomerRoute(admin.Group("/customer"))
	customerinstallation.AdminCustomerInstallationRoute(admin.Group("/customer-installation"))
	asset.AdminAssetRoute(admin.Group("/asset"))
	transaction.AdminTrasactionRoute(admin.Group("/transaction"))
	invoice.AdminInvoiceRoute(admin.Group("/invoice"))

	customer := api.Group("/customer")
	customerdashboard.CustomerDashboardRoute(customer.Group("/dashboard"))

	webhook := api.Group("/webhook")
	moota.WebhookMootaRoute(webhook.Group("/moota"))

	webhookPlain := app.Group("/webhook")
	midtrans.WebhookMidtransRoute(webhookPlain.Group("/midtrans"))

}
