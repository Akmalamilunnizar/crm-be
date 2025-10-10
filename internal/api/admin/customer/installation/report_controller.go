package customerinstallation

import (
	"log"
	"skripsi-be/internal/helpers"
	"strings"

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

// GetCompleteInstallationReportWithTechnicianPhotos - Get complete installation report with computed technician photos
func (c AdminInstallationReportControllerStruct) GetCompleteInstallationReportWithTechnicianPhotos(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	log.Printf("=== CONTROLLER DEBUG: GetCompleteInstallationReportWithTechnicianPhotos ===")
	log.Printf("Installation ID: %s", installationId)

	if installationId == "" {
		log.Printf("ERROR: Installation ID is empty")
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	response, err := c.service.GetCompleteInstallationReportWithTechnicianPhotosService(installationId)
	if err != nil {
		log.Printf("ERROR: Failed to get installation report: %v", err)
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get installation report with technician photos", err.Error())
	}

	log.Printf("SUCCESS: Retrieved installation report with technician photos")
	log.Printf("Response contains - Customer: %s, Technician: %s, Images: %d",
		response.CustomerName, response.TechnicianName, len(response.Images))
	log.Printf("=== END CONTROLLER DEBUG ===")

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation report with technician photos retrieved successfully", response)
}

// GetCompleteInstallationReportByView - Get complete installation report using database view
func (c AdminInstallationReportControllerStruct) GetCompleteInstallationReportByView(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	// Validate installation ID parameter
	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	// Validate UUID format
	if len(installationId) != 36 {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Invalid installation ID format", nil)
	}

	// Get report from service
	report, err := c.service.GetCompleteInstallationReportByViewService(installationId)
	if err != nil {
		log.Printf("Error retrieving installation report for ID %s: %v", installationId, err)

		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") {
			return helpers.ResponseUtils(ctx, fiber.StatusNotFound, false, "Installation report not found", nil)
		}

		// Return generic error for other cases
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to retrieve installation report", nil)
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

// GetInstallationTechnicianTeam - Get technician team for specific installation
func (c AdminInstallationReportControllerStruct) GetInstallationTechnicianTeam(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	technicians, err := c.service.GetInstallationTechnicianTeamService(installationId)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to get technician team", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Technician team retrieved successfully", technicians)
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

// UpdateCompleteInstallationReport - Update complete installation report with all related data
func (c AdminInstallationReportControllerStruct) UpdateCompleteInstallationReport(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")
	log.Printf("DEBUG: UpdateCompleteInstallationReport called with ID: %s", installationId)

	if installationId == "" {
		log.Printf("ERROR: Installation ID is empty")
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	var request UpdateCompleteInstallationReportRequest

	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("ERROR: Failed to parse request body: %v", err)
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Invalid request body", err.Error())
	}

	// Debug request data
	log.Printf("DEBUG: Request parsed successfully")
	log.Printf("DEBUG: Customer ID: %s", request.CustomerId)
	log.Printf("DEBUG: Technician ID: %s", request.TechnicianId)
	log.Printf("DEBUG: Technician photos count: %d", len(request.TechnicianPhotos))
	log.Printf("DEBUG: Technician photos data: %+v", request.TechnicianPhotos)
	log.Printf("DEBUG: Technician photos notes: %s", request.TechnicianPhotosNotes)

	// Validate required fields
	if request.CustomerId == "" {
		log.Printf("ERROR: Customer ID is required but empty")
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Customer ID is required", nil)
	}
	if request.TechnicianId == "" {
		log.Printf("ERROR: Technician ID is required but empty")
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Technician ID is required", nil)
	}

	// Validate technician photos
	if len(request.TechnicianPhotos) > 10 {
		log.Printf("ERROR: Technician photos count exceeds limit: %d > 10", len(request.TechnicianPhotos))
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Maximum 10 technician photos allowed", nil)
	}

	log.Printf("DEBUG: Calling service to update installation report")
	installation, err := c.service.UpdateCompleteInstallationReportService(installationId, request)
	if err != nil {
		log.Printf("ERROR: Service failed to update installation report: %v", err)
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to update installation report", err.Error())
	}

	log.Printf("DEBUG: Installation report updated successfully")
	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation report updated successfully", installation)
}

// DeleteInstallationReport - Delete installation report and update MAC address status
func (c AdminInstallationReportControllerStruct) DeleteInstallationReport(ctx *fiber.Ctx) error {
	installationId := ctx.Params("id")

	if installationId == "" {
		return helpers.ResponseUtils(ctx, fiber.StatusBadRequest, false, "Installation ID is required", nil)
	}

	err := c.service.DeleteInstallationReportService(installationId)
	if err != nil {
		return helpers.ResponseUtils(ctx, fiber.StatusInternalServerError, false, "Failed to delete installation report", err.Error())
	}

	return helpers.ResponseUtils(ctx, fiber.StatusOK, true, "Installation report deleted successfully and MAC address status updated", nil)
}
