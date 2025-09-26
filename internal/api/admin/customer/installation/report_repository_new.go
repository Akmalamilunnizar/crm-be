package customerinstallation

import (
	"log"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type ReportInstallationRepositoryInterface interface {
	CreateReportInstallationRepository(request CreateReportInstallationRequest) (entities.CustomerInstallation, error)
}

type ReportInstallationRepository struct {
	db *gorm.DB
}

func NewReportInstallationRepository(db *gorm.DB) ReportInstallationRepositoryInterface {
	return &ReportInstallationRepository{db: db}
}

// CreateReportInstallationRepository - Create installation report with all related data
func (r *ReportInstallationRepository) CreateReportInstallationRepository(request CreateReportInstallationRequest) (entities.CustomerInstallation, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Parse dates
	var onAirDate, trialEndDate, serviceReadyDate, installationCompletedAt *time.Time

	if request.OnAirDate != "" {
		if parsed, err := time.Parse("2006-01-02", request.OnAirDate); err == nil {
			onAirDate = &parsed
		}
	}

	if request.TrialEndDate != "" {
		if parsed, err := time.Parse("2006-01-02", request.TrialEndDate); err == nil {
			trialEndDate = &parsed
		}
	}

	if request.ServiceReadyDate != "" {
		if parsed, err := time.Parse("2006-01-02", request.ServiceReadyDate); err == nil {
			serviceReadyDate = &parsed
		}
	}

	if request.InstallationCompletedAt != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", request.InstallationCompletedAt); err == nil {
			installationCompletedAt = &parsed
		}
	}

	// Set default values
	if request.Status == "" {
		request.Status = "pending"
	}
	if request.InstallationType == "" {
		request.InstallationType = "new_installation"
	}
	if request.UserStatus == "" {
		request.UserStatus = "Active"
	}
	if request.StatusPerangkat == "" {
		request.StatusPerangkat = "active"
	}
	if request.KepemilikanPerangkat == "" {
		request.KepemilikanPerangkat = "owned"
	}
	if request.LastPingStatus == "" {
		request.LastPingStatus = "unknown"
	}

	// Create main installation record
	installation := entities.CustomerInstallation{
		CustomerID:              &request.CustomerID,
		TechnicianID:            &request.TechnicianID,
		Status:                  request.Status,
		Notes:                   request.Notes,
		DocumentType:            &request.DocumentType,
		DocumentPhoto:           &request.DocumentPhoto,
		InstallationType:        request.InstallationType,
		TotalAssetsOut:          0, // Will be calculated based on asset transactions
		TotalAssetsIn:           0, // Will be calculated based on asset transactions
		InstallationCompletedAt: installationCompletedAt,
		TrialEndDate:            trialEndDate,
		ServiceReadyDate:        serviceReadyDate,
		OnAirDate:               onAirDate,
	}

	// Log installation data before creating
	var documentPhotoStr, documentTypeStr string
	if installation.DocumentPhoto != nil {
		documentPhotoStr = *installation.DocumentPhoto
	} else {
		documentPhotoStr = "nil"
	}
	if installation.DocumentType != nil {
		documentTypeStr = *installation.DocumentType
	} else {
		documentTypeStr = "nil"
	}
	log.Printf("Creating installation record - DocumentPhoto: %s, DocumentType: %s",
		documentPhotoStr, documentTypeStr)

	if err := tx.Create(&installation).Error; err != nil {
		log.Printf("Failed to create installation record: %v", err)
		tx.Rollback()
		return installation, err
	}

	log.Printf("Installation record created successfully with ID: %s", installation.ID)

	// Create network device
	if request.AssetsID != "" {
		networkDevice := entities.NetworkDevice{
			ID:                     "",
			CustomerID:             request.CustomerID,
			AssetsID:               request.AssetsID,
			CustomerInstallationID: &installation.ID,
			SwitchID:               &request.SwitchID,
			PortNumber:             &request.PortNumber,
			RemotePort:             &request.RemotePort,
			EthPort:                &request.EthPort,
			MacAddress:             &request.MacAddress,
			IPStatic:               &request.IPStatic,
			KepemilikanPerangkat:   request.KepemilikanPerangkat,
			StatusPerangkat:        request.StatusPerangkat,
			LastPingStatus:         request.LastPingStatus,
		}

		if err := tx.Create(&networkDevice).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create customer service
	if request.UserLogin != "" || request.Password != "" {
		customerService := entities.CustomerService{
			ID:                     "",
			CustomerID:             request.CustomerID,
			CustomerInstallationID: &installation.ID,
			UserLogin:              &request.UserLogin,
			Password:               &request.Password,
			UserStatus:             request.UserStatus,
			EndPortType:            &request.EndPortType,
			InstallationNotes:      &request.InstallationNotes,
		}

		if err := tx.Create(&customerService).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create cable
	if request.CableType != "" || request.CableLength > 0 {
		cable := entities.Cable{
			ID:                     "",
			Name:                   "Installation Cable",
			Type:                   &request.CableType,
			Length:                 &request.CableLength,
			Status:                 "in_use",
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&cable).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return installation, err
	}

	// Load the created installation with relations
	var result entities.CustomerInstallation
	err := r.db.Preload("Customer").
		Preload("Technician").
		Preload("NetworkDevices").
		Preload("NetworkDevices.Assets").
		Preload("CustomerServices").
		Preload("Cables").
		Where("id = ?", installation.ID).
		First(&result).Error

	return result, err
}
