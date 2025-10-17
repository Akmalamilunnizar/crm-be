package customerinstallation

import (
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"skripsi-be/internal/models/entities"
)

type AdminCustomerInstallationRepositoryInterface interface {
	CreateAdminCustomerInstallationRepository(customer CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	FindByIdAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	FindAdminCustomerInstallationRepository() ([]entities.CustomerInstallation, error)
	GetInstallationReportsByCustomerRepository(customerId string) ([]entities.CustomerInstallation, error) // NEW
}
type AdminCustomerInstallationRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminCustomerInstallationRepository(db *gorm.DB) AdminCustomerInstallationRepositoryStruct {
	return AdminCustomerInstallationRepositoryStruct{db}
}

func (r AdminCustomerInstallationRepositoryStruct) FindAdminCustomerInstallationRepository() ([]entities.CustomerInstallation, error) {
	customerInstallations := []entities.CustomerInstallation{}
	tx := r.db.Preload("Customer", "deleted_at IS NULL").Preload("Technician").Preload("Images").Find(&customerInstallations)

	if tx.Error != nil {
		return customerInstallations, tx.Error
	}

	return customerInstallations, nil
}
func (r AdminCustomerInstallationRepositoryStruct) FindByIdAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Customer", "deleted_at IS NULL").Preload("Technician").Preload("Images").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}
func (r AdminCustomerInstallationRepositoryStruct) CreateAdminCustomerInstallationRepository(request CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	images := []entities.Image{}

	// Parse on_air_date if provided
	var onAirDate *time.Time
	if request.OnAirDate != "" {
		parsedOnAirDate, err := time.Parse("2006-01-02", request.OnAirDate)
		if err != nil {
			return entities.CustomerInstallation{}, err
		}
		onAirDate = &parsedOnAirDate
	}

	// Handle document_type - convert to pointer if not empty
	var documentType *string
	if request.DocumentType != "" {
		documentType = &request.DocumentType
	}

	// Handle document_photo - convert to pointer if not empty
	var documentPhoto *string
	if request.DocumentPhoto != "" {
		documentPhoto = &request.DocumentPhoto
	}

	customerInstallation := entities.CustomerInstallation{
		ID:            "",
		CustomerID:    &request.CustomerId,
		TechnicianID:  &request.TechnicianId,
		Status:        request.Status,
		Notes:         request.Notes,
		DocumentType:  documentType,
		DocumentPhoto: documentPhoto,
		OnAirDate:     onAirDate,
	}

	tx := r.db.Begin()
	txCustomer := tx.Create(&customerInstallation)
	if txCustomer.Error != nil {
		tx.Rollback()
		return entities.CustomerInstallation{}, txCustomer.Error
	}

	// Update images with the installation ID
	if len(request.ImageIds) > 0 {
		txImage := tx.Where("id IN ?", request.ImageIds).Find(&images).Update("archive_installation_id", customerInstallation.ID)
		if txImage.Error != nil {
			tx.Rollback()
			return entities.CustomerInstallation{}, txImage.Error
		}
	}

	if tx.Error != nil {
		tx.Rollback()
		return entities.CustomerInstallation{}, tx.Error
	}

	tx.Commit()
	return customerInstallation, nil

}
func (r AdminCustomerInstallationRepositoryStruct) UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Save(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}

func (r AdminCustomerInstallationRepositoryStruct) DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Delete(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}

// NEW: Repository method to get installation reports by customer ID
func (r AdminCustomerInstallationRepositoryStruct) GetInstallationReportsByCustomerRepository(customerId string) ([]entities.CustomerInstallation, error) {
	reports := []entities.CustomerInstallation{}

	// Get all installation reports for the specified customer with related data
	tx := r.db.Preload("Customer", "deleted_at IS NULL").Preload("Technician").Preload("Images").
		Where("customer_id = ?", customerId).
		Order("created_at DESC").
		Find(&reports)

	if tx.Error != nil {
		return reports, tx.Error
	}

	return reports, nil
}
