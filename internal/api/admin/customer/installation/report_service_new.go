package customerinstallation

import (
	"fmt"
	"log"
	"skripsi-be/internal/api/admin/recurring_invoice"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type ReportInstallationServiceInterface interface {
	CreateReportInstallationService(request CreateReportInstallationRequest, createdByUserID string) (entities.CustomerInstallation, error)
	DeleteInstallation(installationId string) error
	UpdateInstallationCodeName(installationId string, codeName string) error
	GetDB() *gorm.DB
}

type ReportInstallationService struct {
	repository           ReportInstallationRepositoryInterface
	recurringInvoiceRepo recurring_invoice.AdminRecurringInvoiceRepositoryInterface
}

func NewReportInstallationService(repository ReportInstallationRepositoryInterface, recurringInvoiceRepo recurring_invoice.AdminRecurringInvoiceRepositoryInterface) ReportInstallationServiceInterface {
	return &ReportInstallationService{
		repository:           repository,
		recurringInvoiceRepo: recurringInvoiceRepo,
	}
}

// CreateReportInstallationService - Create installation report with all related data
func (s *ReportInstallationService) CreateReportInstallationService(request CreateReportInstallationRequest, createdByUserID string) (entities.CustomerInstallation, error) {
	installation, err := s.repository.CreateReportInstallationRepository(request, createdByUserID)
	if err != nil {
		return installation, err
	}

	// Create automatic recurring invoice for the customer
	err = s.createRecurringInvoiceForInstallation(installation, request, createdByUserID)
	if err != nil {
		log.Printf("Warning: Failed to create recurring invoice for installation %s: %v", installation.ID, err)
		// Don't fail the installation creation if recurring invoice fails
	}

	return installation, nil
}

// DeleteInstallation - Delete installation by ID
func (s *ReportInstallationService) DeleteInstallation(installationId string) error {
	return s.repository.DeleteInstallation(installationId)
}

// createRecurringInvoiceForInstallation - Create a recurring invoice for the installation
func (s *ReportInstallationService) createRecurringInvoiceForInstallation(installation entities.CustomerInstallation, request CreateReportInstallationRequest, createdByUserID string) error {
	log.Printf("🔄 Creating recurring invoice for installation: %s", installation.ID)

	// Get product information from the request
	var product entities.Products
	if request.ProductID != "" {
		// Get product details from database
		db := s.repository.GetDB() // We need to add this method to the repository interface
		if db != nil {
			if err := db.Where("id = ?", request.ProductID).First(&product).Error; err != nil {
				log.Printf("❌ Failed to get product details: %v", err)
				return err
			}
		}
	}

	// Set default values if product not found
	if product.ID == "" {
		product.Name = "Internet Service"
		product.Price = 100000 // Default 100k IDR
		product.ID = "default"
	}

	log.Printf("📦 Product info: Name=%s, Price=%d", product.Name, product.Price)

	// Create invoice item name with product details
	itemName := fmt.Sprintf("%s - %s", product.Name, installation.ID[:8])

	// Create invoice items
	items := []recurring_invoice.RecurringInvoiceItem{
		{
			Name:  itemName,
			Price: product.Price,
			Qty:   1,
			Total: product.Price,
		},
	}

	// Set up dates - use installation_completed_at if available, otherwise current time
	var invoiceDate, dueDate time.Time

	if installation.InstallationCompletedAt != nil {
		invoiceDate = *installation.InstallationCompletedAt
		// Due date is one month from invoice date
		dueDate = invoiceDate.AddDate(0, 1, 0)
		log.Printf("📅 Using InstallationCompletedAt: invoice_date=%s, due_date=%s",
			invoiceDate.Format("2006-01-02"), dueDate.Format("2006-01-02"))
	} else {
		// Fallback to current time
		invoiceDate = time.Now()
		dueDate = invoiceDate.AddDate(0, 1, 0)
		log.Printf("📅 Using current time: invoice_date=%s, due_date=%s",
			invoiceDate.Format("2006-01-02"), dueDate.Format("2006-01-02"))
	}

	// Create recurring invoice request
	recurringRequest := recurring_invoice.CreateRecurringInvoiceRequest{
		CustomerID:           *installation.CustomerID,
		CustomerInstallation: &installation.ID,
		Amount:               product.Price,
		InvoiceDate:          invoiceDate,
		DueDate:              dueDate,
		Frequency:            "monthly",
		Description:          stringPtrHelper(fmt.Sprintf("Monthly internet service for %s", product.Name)),
		InvoiceItems:         items,
	}

	// Create the recurring invoice
	_, err := s.recurringInvoiceRepo.CreateRecurringInvoice(recurringRequest, createdByUserID)
	if err != nil {
		return err
	}

	log.Printf("✅ Successfully created recurring invoice for customer %s (installation %s)", *installation.CustomerID, installation.ID)
	return nil
}

// UpdateInstallationCodeName updates the code_name field for an installation
func (s *ReportInstallationService) UpdateInstallationCodeName(installationId string, codeName string) error {
	return s.repository.UpdateInstallationCodeName(installationId, codeName)
}

// GetDB returns the database connection for direct queries
func (s *ReportInstallationService) GetDB() *gorm.DB {
	return s.repository.GetDB()
}

// Helper function to create string pointer
func stringPtrHelper(s string) *string {
	return &s
}
