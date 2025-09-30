package customerinstallation

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type AdminInstallationReportControllerStruct struct {
	service AdminInstallationReportServiceInterface
}

func NewAdminInstallationReportController(service AdminInstallationReportServiceInterface) AdminInstallationReportControllerStruct {
	return AdminInstallationReportControllerStruct{service}
}

// GetCompleteInstallationReport - Get complete installation report with all related data
func (c AdminInstallationReportControllerStruct) GetCompleteInstallationReport(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	installation, err := c.service.GetCompleteInstallationReportService(installationId)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get installation report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation report retrieved successfully", installation)
}

// GetCompleteInstallationReportByView - Get complete installation report using database view
func (c AdminInstallationReportControllerStruct) GetCompleteInstallationReportByView(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	report, err := c.service.GetCompleteInstallationReportByViewService(installationId)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get installation report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation report retrieved successfully", report)
}

// GetAllCompleteInstallationReports - Get all complete installation reports
func (c AdminInstallationReportControllerStruct) GetAllCompleteInstallationReports(ctx *fiber.Ctx) error {
	reports, err := c.service.GetAllCompleteInstallationReportsService()
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get installation reports", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation reports retrieved successfully", reports)
}

// GetInstallationSummaryPerCustomer - Get installation summary grouped by customer
func (c AdminInstallationReportControllerStruct) GetInstallationSummaryPerCustomer(ctx *fiber.Ctx) error {
	summaries, err := c.service.GetInstallationSummaryPerCustomerService()
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get installation summary", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation summary retrieved successfully", summaries)
}

// GetInstallationAssetReport - Get asset report for specific installation
func (c AdminInstallationReportControllerStruct) GetInstallationAssetReport(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	report, err := c.service.GetInstallationAssetReportService(installationId)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get asset report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Asset report retrieved successfully", report)
}

// GetInstallationTechnicianReport - Get technician performance report
func (c AdminInstallationReportControllerStruct) GetInstallationTechnicianReport(ctx *fiber.Ctx) error {
	reports, err := c.service.GetInstallationTechnicianReportService()
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get technician report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Technician report retrieved successfully", reports)
}

// CreateCompleteInstallationReport - Create complete installation report with all related data
func (c AdminInstallationReportControllerStruct) CreateCompleteInstallationReport(ctx *fiber.Ctx) error {
	var request CreateCompleteInstallationReportRequest

	if err := ctx.BodyParser(&request); err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Invalid request body", err.Error())
	}

	// Validate required fields
	if request.CustomerId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Customer ID is required", nil)
	}
	if request.TechnicianId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Technician ID is required", nil)
	}
	if len(request.ImageIds) == 0 {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "At least one image is required", nil)
	}

	installation, err := c.service.CreateCompleteInstallationReportService(request)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to create installation report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusCreated, true, "Installation report created successfully", installation)
}
