package customerinstallation

import (
	"fmt"
	"log"
	"skripsi-be/internal/api/admin/recurring_invoice"
	"skripsi-be/internal/models/entities"
	"time"
)

type AdminInstallationReportServiceInterface interface {
	GetCompleteInstallationReportService(installationId string) (entities.CustomerInstallation, error)
	GetCompleteInstallationReportWithTechnicianPhotosService(installationId string) (CompleteInstallationReportWithTechnicianPhotosResponse, error)
	GetCompleteInstallationReportByViewService(installationId string) (InstallationReportCompleteResponse, error)
	GetAllCompleteInstallationReportsService() ([]InstallationReportCompleteResponse, error)
	GetInstallationSummaryPerCustomerService() ([]InstallationSummaryResponse, error)
	GetInstallationAssetReportService(installationId string) (InstallationAssetReportResponse, error)
	GetInstallationTechnicianReportService() ([]InstallationTechnicianReportResponse, error)
	GetInstallationTechnicianTeamService(installationId string) ([]InstallationTechnicianTeamResponse, error)
	CreateCompleteInstallationReportService(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
	UpdateCompleteInstallationReportService(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
	UpdateInstallationCodeName(installationId string, codeName string) error
	DeleteInstallationReportService(installationId string) error
}

type AdminInstallationReportServiceStruct struct {
	repository           AdminInstallationReportRepositoryInterface
	recurringInvoiceRepo recurring_invoice.AdminRecurringInvoiceRepositoryInterface
}

func NewAdminInstallationReportService(repository AdminInstallationReportRepositoryInterface, recurringInvoiceRepo recurring_invoice.AdminRecurringInvoiceRepositoryInterface) AdminInstallationReportServiceStruct {
	return AdminInstallationReportServiceStruct{
		repository:           repository,
		recurringInvoiceRepo: recurringInvoiceRepo,
	}
}

// GetCompleteInstallationReportService - Get complete installation report with all related data
func (s AdminInstallationReportServiceStruct) GetCompleteInstallationReportService(installationId string) (entities.CustomerInstallation, error) {
	installation, err := s.repository.FindCompleteInstallationReportRepository(installationId)
	if err != nil {
		return installation, err
	}
	return installation, nil
}

// GetCompleteInstallationReportWithTechnicianPhotosService - Get complete installation report with computed technician photos
func (s AdminInstallationReportServiceStruct) GetCompleteInstallationReportWithTechnicianPhotosService(installationId string) (CompleteInstallationReportWithTechnicianPhotosResponse, error) {
	response, err := s.repository.FindCompleteInstallationReportWithTechnicianPhotosRepository(installationId)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetCompleteInstallationReportByViewService - Get complete installation report using database view
func (s AdminInstallationReportServiceStruct) GetCompleteInstallationReportByViewService(installationId string) (InstallationReportCompleteResponse, error) {
	// Validate input
	if installationId == "" {
		return InstallationReportCompleteResponse{}, fmt.Errorf("installation ID cannot be empty")
	}

	// Validate UUID format
	if len(installationId) != 36 {
		return InstallationReportCompleteResponse{}, fmt.Errorf("invalid installation ID format")
	}

	report, err := s.repository.FindCompleteInstallationReportByViewRepository(installationId)
	if err != nil {
		log.Printf("Error retrieving installation report for ID %s: %v", installationId, err)
		return InstallationReportCompleteResponse{}, fmt.Errorf("failed to retrieve installation report: %w", err)
	}

	// Validate that report was found
	if report.InstallationId == "" {
		return InstallationReportCompleteResponse{}, fmt.Errorf("installation report not found for ID: %s", installationId)
	}

	return report, nil
}

// GetAllCompleteInstallationReportsService - Get all complete installation reports
func (s AdminInstallationReportServiceStruct) GetAllCompleteInstallationReportsService() ([]InstallationReportCompleteResponse, error) {
	reports, err := s.repository.FindAllCompleteInstallationReportsRepository()
	if err != nil {
		return reports, err
	}
	return reports, nil
}

// GetInstallationSummaryPerCustomerService - Get installation summary grouped by customer
func (s AdminInstallationReportServiceStruct) GetInstallationSummaryPerCustomerService() ([]InstallationSummaryResponse, error) {
	summaries, err := s.repository.FindInstallationSummaryPerCustomerRepository()
	if err != nil {
		return summaries, err
	}
	return summaries, nil
}

// GetInstallationAssetReportService - Get asset report for specific installation
func (s AdminInstallationReportServiceStruct) GetInstallationAssetReportService(installationId string) (InstallationAssetReportResponse, error) {
	report, err := s.repository.FindInstallationAssetReportRepository(installationId)
	if err != nil {
		return report, err
	}
	return report, nil
}

// GetInstallationTechnicianReportService - Get technician performance report
func (s AdminInstallationReportServiceStruct) GetInstallationTechnicianReportService() ([]InstallationTechnicianReportResponse, error) {
	reports, err := s.repository.FindInstallationTechnicianReportRepository()
	if err != nil {
		return reports, err
	}
	return reports, nil
}

// GetInstallationTechnicianTeamService - Get technician team for specific installation
func (s AdminInstallationReportServiceStruct) GetInstallationTechnicianTeamService(installationId string) ([]InstallationTechnicianTeamResponse, error) {
	technicians, err := s.repository.FindInstallationTechnicianTeamRepository(installationId)
	if err != nil {
		return technicians, err
	}
	return technicians, nil
}

// CreateCompleteInstallationReportService - Create complete installation report with all related data
func (s AdminInstallationReportServiceStruct) CreateCompleteInstallationReportService(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
	installation, err := s.repository.CreateCompleteInstallationReportRepository(request)
	if err != nil {
		return installation, err
	}

	// Create automatic recurring invoice for the customer
	err = s.createRecurringInvoiceForInstallation(installation, request)
	if err != nil {
		log.Printf("Warning: Failed to create recurring invoice for installation %s: %v", installation.ID, err)
		// Don't fail the installation creation if recurring invoice fails
	}

	return installation, nil
}

// createRecurringInvoiceForInstallation - Create a recurring invoice for the installation
func (s AdminInstallationReportServiceStruct) createRecurringInvoiceForInstallation(installation entities.CustomerInstallation, request CreateCompleteInstallationReportRequest) error {
	// Get customer information to determine package pricing
	// For now, we'll use a default monthly internet service fee
	// In a real implementation, you'd fetch the customer's package details from the products table

	defaultPrice := int64(100000) // Default 100k IDR per month
	serviceName := "Internet Service"

	// Create invoice items
	items := []recurring_invoice.RecurringInvoiceItem{
		{
			Name:  serviceName,
			Price: defaultPrice,
			Qty:   1,
			Total: defaultPrice,
		},
	}

	// Set up dates
	now := time.Now()
	invoiceDate := now
	dueDate := now.AddDate(0, 1, 0) // One month from now

	// Create recurring invoice request
	recurringRequest := recurring_invoice.CreateRecurringInvoiceRequest{
		CustomerID:           *installation.CustomerID,
		CustomerInstallation: &installation.ID,
		Amount:               defaultPrice,
		InvoiceDate:          invoiceDate,
		DueDate:              dueDate,
		Frequency:            "monthly",
		Description:          stringPtr("Automatic recurring invoice for internet service"),
		InvoiceItems:         items,
	}

	// Create the recurring invoice
	_, err := s.recurringInvoiceRepo.CreateRecurringInvoice(recurringRequest, "system")
	if err != nil {
		return err
	}

	log.Printf("Successfully created recurring invoice for customer %s (installation %s)", *installation.CustomerID, installation.ID)
	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// UpdateCompleteInstallationReportService - Update complete installation report with all related data
func (s AdminInstallationReportServiceStruct) UpdateCompleteInstallationReportService(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
	installation, err := s.repository.UpdateCompleteInstallationReportRepository(installationId, request)
	if err != nil {
		return installation, err
	}
	return installation, nil
}

// UpdateInstallationCodeName updates the code_name field for an installation
func (s AdminInstallationReportServiceStruct) UpdateInstallationCodeName(installationId string, codeName string) error {
	return s.repository.UpdateInstallationCodeName(installationId, codeName)
}

// DeleteInstallationReportService - Delete installation report and update MAC address status
func (s AdminInstallationReportServiceStruct) DeleteInstallationReportService(installationId string) error {
	err := s.repository.DeleteInstallationReportRepository(installationId)
	if err != nil {
		return err
	}
	return nil
}
