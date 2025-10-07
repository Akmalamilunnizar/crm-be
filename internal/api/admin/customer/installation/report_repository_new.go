package customerinstallation

import (
	"database/sql"
	"fmt"
	"log"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type ReportInstallationRepositoryInterface interface {
	CreateReportInstallationRepository(request CreateReportInstallationRequest, createdByUserID string) (entities.CustomerInstallation, error)
	CreateInstallationTechnicians(tx *gorm.DB, installationID string, technicians []TechnicianAssignment) error
	DeleteInstallation(installationId string) error
}

type ReportInstallationRepository struct {
	db *gorm.DB
}

func NewReportInstallationRepository(db *gorm.DB) ReportInstallationRepositoryInterface {
	return &ReportInstallationRepository{db: db}
}

// CreateInstallationTechnicians creates multiple technician assignments for an installation
func (r *ReportInstallationRepository) CreateInstallationTechnicians(
	tx *gorm.DB,
	installationID string,
	technicians []TechnicianAssignment,
) error {
	if len(technicians) == 0 {
		return nil
	}

	// Validate at least one senior exists
	hasSenior := false
	for _, tech := range technicians {
		if tech.Role == "senior" {
			hasSenior = true
			break
		}
	}

	if !hasSenior {
		return fmt.Errorf("at least one senior technician is required")
	}

	// Ensure only one primary technician
	primaryCount := 0
	for i := range technicians {
		if technicians[i].IsPrimary {
			primaryCount++
		}
	}

	// If no primary is set, make the first senior primary
	if primaryCount == 0 {
		for i := range technicians {
			if technicians[i].Role == "senior" {
				technicians[i].IsPrimary = true
				break
			}
		}
	}

	// Create pivot records
	for _, tech := range technicians {
		pivot := entities.InstallationReportTechnician{
			CustomerInstallationID: installationID,
			TechnicianID:           tech.TechnicianID,
			Role:                   tech.Role,
			IsPrimary:              tech.IsPrimary,
			Notes:                  tech.Notes,
		}
		if err := tx.Create(&pivot).Error; err != nil {
			return fmt.Errorf("failed to create technician assignment: %w", err)
		}
	}

	log.Printf("✅ Created %d technician assignments for installation %s", len(technicians), installationID)
	return nil
}

// Helper function to convert empty string to nil
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// CreateReportInstallationRepository - Create installation report with all related data
func (r *ReportInstallationRepository) CreateReportInstallationRepository(request CreateReportInstallationRequest, createdByUserID string) (entities.CustomerInstallation, error) {
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
		// Try different date formats
		formats := []string{
			"2006-01-02T15:04",    // datetime-local format from frontend
			"2006-01-02 15:04:05", // standard datetime format
			"2006-01-02T15:04:05", // ISO format
		}

		for _, format := range formats {
			if parsed, err := time.Parse(format, request.InstallationCompletedAt); err == nil {
				installationCompletedAt = &parsed
				break
			}
		}

		// Log for debugging
		if installationCompletedAt != nil {
			log.Printf("Successfully parsed installation_completed_at: %s -> %v", request.InstallationCompletedAt, *installationCompletedAt)
		} else {
			log.Printf("Failed to parse installation_completed_at: %s", request.InstallationCompletedAt)
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

	// Determine technician_id for backward compatibility
	var technicianID *string
	if request.TechnicianID != "" {
		// Legacy: single technician provided
		technicianID = &request.TechnicianID
	} else if len(request.Technicians) > 0 {
		// New: multiple technicians - set primary technician as the legacy technician_id
		for _, tech := range request.Technicians {
			if tech.IsPrimary {
				technicianID = &tech.TechnicianID
				break
			}
		}
		// If no primary found, use first senior technician
		if technicianID == nil {
			for _, tech := range request.Technicians {
				if tech.Role == "senior" {
					technicianID = &tech.TechnicianID
					break
				}
			}
		}
	}

	// Create main installation record
	installation := entities.CustomerInstallation{
		CustomerID:              stringToPtr(request.CustomerID),
		TechnicianID:            technicianID,
		Status:                  request.Status,
		Notes:                   request.Notes,
		DocumentType:            stringToPtr(request.DocumentType),
		DocumentPhoto:           stringToPtr(request.DocumentPhoto),
		InstallationType:        request.InstallationType,
		TotalAssetsOut:          0, // Will be calculated based on asset transactions
		TotalAssetsIn:           0, // Will be calculated based on asset transactions
		InstallationCompletedAt: installationCompletedAt,
		TrialEndDate:            trialEndDate,
		ServiceReadyDate:        serviceReadyDate,
		OnAirDate:               onAirDate,
	}

	// Normalize document photo path before saving
	if request.DocumentPhoto != "" {
		request.DocumentPhoto = normalizeDocumentPhotoPath(request.DocumentPhoto)
		log.Printf("Normalized document photo path from '%s' to '%s'", request.DocumentPhoto, request.DocumentPhoto)
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
	log.Printf("=== REPOSITORY DEBUG ===")
	log.Printf("Creating installation record - DocumentPhoto: %s, DocumentType: %s",
		documentPhotoStr, documentTypeStr)
	log.Printf("Request DocumentPhoto field: '%s'", request.DocumentPhoto)
	log.Printf("Request DocumentType field: '%s'", request.DocumentType)
	log.Printf("=== END REPOSITORY DEBUG ===")

	if err := tx.Create(&installation).Error; err != nil {
		log.Printf("Failed to create installation record: %v", err)
		tx.Rollback()
		return installation, err
	}

	log.Printf("Installation record created successfully with ID: %s", installation.ID)

	// Log what was actually saved to database
	log.Printf("=== DATABASE SAVE DEBUG ===")
	if installation.DocumentPhoto != nil {
		log.Printf("✅ Document Photo saved to DB: '%s'", *installation.DocumentPhoto)
	} else {
		log.Printf("❌ Document Photo is NULL in DB")
	}
	if installation.DocumentType != nil {
		log.Printf("✅ Document Type saved to DB: '%s'", *installation.DocumentType)
	} else {
		log.Printf("❌ Document Type is NULL in DB")
	}
	log.Printf("=== END DATABASE SAVE DEBUG ===")

	// Update customer with company_id and sales_representative_id if provided
	if request.CustomerCompanyID != "" || request.CustomerSalesRepresentativeID != "" {
		updateData := make(map[string]interface{})
		if request.CustomerCompanyID != "" {
			updateData["company_id"] = request.CustomerCompanyID
		}
		if request.CustomerSalesRepresentativeID != "" {
			updateData["sales_representative_id"] = request.CustomerSalesRepresentativeID
		}

		if err := tx.Model(&entities.Customer{}).
			Where("id = ?", request.CustomerID).
			Updates(updateData).Error; err != nil {
			log.Printf("Failed to update customer with company/sales rep info: %v", err)
			tx.Rollback()
			return installation, err
		}
		log.Printf("Updated customer with company_id: %s, sales_representative_id: %s",
			request.CustomerCompanyID, request.CustomerSalesRepresentativeID)
	}

	// Create network device - only if essential fields are provided
	if request.AssetsID != "" && (request.MacAddress != "" || request.IPStatic != "" || request.SwitchID != "") {
		// Prepare ProductID - only set if not empty
		var productID *string
		if request.ProductID != "" {
			productID = &request.ProductID
		}

		networkDevice := entities.NetworkDevice{
			ID:         "",
			CustomerID: request.CustomerID,
			AssetsID: sql.NullString{
				String: request.AssetsID,
				Valid:  request.AssetsID != "",
			},
			CustomerInstallationID: &installation.ID,
			SwitchID:               stringToPtr(request.SwitchID),
			PortNumber:             stringToPtr(request.PortNumber),
			RemotePort:             stringToPtr(request.RemotePort),
			EthPort:                stringToPtr(request.EthPort),
			MacAddress:             stringToPtr(request.MacAddress),
			IPStatic:               stringToPtr(request.IPStatic),
			KepemilikanPerangkat:   request.KepemilikanPerangkat,
			StatusPerangkat:        request.StatusPerangkat,
			LastPingStatus:         request.LastPingStatus,
			ProductID:              productID,
		}

		if err := tx.Create(&networkDevice).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
		log.Printf("Created network device for installation %s", installation.ID)

		// Create asset transaction to track asset item assignment
		log.Printf("=== ASSET TRANSACTION DEBUG ===")
		log.Printf("Request.MacAddress: '%s'", request.MacAddress)
		log.Printf("Request.AssetItemID: '%s'", request.AssetItemID)
		log.Printf("Condition check: MacAddress != '' (%t) || AssetItemID != '' (%t) = %t",
			request.MacAddress != "", request.AssetItemID != "",
			request.MacAddress != "" || request.AssetItemID != "")

		if request.MacAddress != "" || request.AssetItemID != "" {
			var assetItem entities.AssetItem
			var err error

			// Try to find asset item by asset_item_id first, then by MAC address
			if request.AssetItemID != "" {
				err = tx.Where("id = ? AND status = 'in_stock'", request.AssetItemID).First(&assetItem).Error
				if err == nil {
					log.Printf("Found asset item by ID: %s", request.AssetItemID)
				}
			}

			// If not found by ID, try by MAC address
			if err != nil && request.MacAddress != "" {
				err = tx.Where("mac_address = ? AND status = 'in_stock'", request.MacAddress).First(&assetItem).Error
				if err == nil {
					log.Printf("Found asset item by MAC address: %s", request.MacAddress)
				}
			}

			if err == nil {
				// Create asset transaction record
				assetTransaction := entities.AssetTransaction{
					ID:                     "",
					CustomerInstallationID: installation.ID,
					AssetItemID:            &assetItem.ID, // Link to specific asset item
					AssetID:                assetItem.AssetID,
					TransactionType:        "out", // Asset going out to customer
					Quantity:               1,
					Notes:                  stringToPtr(fmt.Sprintf("Asset assigned to installation - MAC: %s, Item ID: %s", assetItem.MacAddress, assetItem.ID)),
					TransactionDate:        time.Now(),
					CreatedBy:              createdByUserID, // Use the authenticated user ID
				}

				if err := tx.Create(&assetTransaction).Error; err != nil {
					log.Printf("Failed to create asset transaction: %v", err)
					// Don't rollback - this is not critical enough to fail the entire installation
				} else {
					// Update asset item status to 'in_use'
					err = tx.Model(&assetItem).Update("status", "in_use").Error
					if err != nil {
						log.Printf("Failed to update asset item status: %v", err)
					} else {
						log.Printf("✅ Created asset transaction and updated asset item %s (MAC: %s) status to 'in_use' for installation %s", assetItem.ID, assetItem.MacAddress, installation.ID)
					}
				}
			} else {
				log.Printf("❌ Asset item not found - AssetItemID: %s, MacAddress: %s, Error: %v", request.AssetItemID, request.MacAddress, err)
			}
		}
	} else {
		log.Printf("Skipping network device creation - insufficient network data provided")
	}

	// Create customer service - only if both login and password are provided
	if request.UserLogin != "" && request.Password != "" {
		customerService := entities.CustomerService{
			ID:                     "",
			CustomerID:             request.CustomerID,
			CustomerInstallationID: &installation.ID,
			UserLogin:              stringToPtr(request.UserLogin),
			Password:               stringToPtr(request.Password),
			UserStatus:             request.UserStatus,
			EndPortType:            stringToPtr(request.EndPortType),
			InstallationNotes:      stringToPtr(request.InstallationNotes),
		}

		if err := tx.Create(&customerService).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
		log.Printf("Created customer service for installation %s", installation.ID)
	} else {
		log.Printf("Skipping customer service creation - login and/or password not provided")
	}

	// Create cable - only if cable type and length are provided
	if request.CableType != "" && request.CableLength > 0 {
		cable := entities.Cable{
			ID:                     "",
			Name:                   "Installation Cable",
			Type:                   stringToPtr(request.CableType),
			Length:                 &request.CableLength,
			Status:                 "in_use",
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&cable).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
		log.Printf("Created cable for installation %s", installation.ID)
	} else {
		log.Printf("Skipping cable creation - cable type and/or length not provided")
	}

	// Update images with installation ID
	if len(request.ImageIds) > 0 {
		log.Printf("Updating %d images with installation ID: %s", len(request.ImageIds), installation.ID)
		if err := tx.Model(&entities.Image{}).
			Where("id IN ?", request.ImageIds).
			Update("archive_installation_id", installation.ID).Error; err != nil {
			log.Printf("Failed to update images with installation ID: %v", err)
			tx.Rollback()
			return installation, err
		}
		log.Printf("Successfully updated images with installation ID")
	}

	// Create technician assignments
	if len(request.Technicians) > 0 {
		log.Printf("Creating technician assignments for installation: %s", installation.ID)
		if err := r.CreateInstallationTechnicians(tx, installation.ID, request.Technicians); err != nil {
			log.Printf("Failed to create technician assignments: %v", err)
			tx.Rollback()
			return installation, err
		}
		log.Printf("Successfully created technician assignments")
	} else if request.TechnicianID != "" {
		// Backward compatibility: if old TechnicianID is provided, create a single senior assignment
		log.Printf("Creating legacy technician assignment for installation: %s", installation.ID)
		legacyTech := []TechnicianAssignment{
			{
				TechnicianID: request.TechnicianID,
				Role:         "senior",
				IsPrimary:    true,
				Notes:        "",
			},
		}
		if err := r.CreateInstallationTechnicians(tx, installation.ID, legacyTech); err != nil {
			log.Printf("Failed to create legacy technician assignment: %v", err)
			tx.Rollback()
			return installation, err
		}
		log.Printf("Successfully created legacy technician assignment")
	}

	if err := tx.Commit().Error; err != nil {
		return installation, err
	}

	// Load the created installation with relations
	var result entities.CustomerInstallation
	err := r.db.Preload("Customer").
		Preload("Technician").
		Preload("InstallationTechnicians").
		Preload("InstallationTechnicians.Technician").
		Preload("Images").
		Preload("NetworkDevices").
		Preload("NetworkDevices.Assets").
		Preload("CustomerServices").
		Preload("Cables").
		Where("id = ?", installation.ID).
		First(&result).Error

	return result, err
}

// DeleteInstallation - Delete installation by ID
func (r *ReportInstallationRepository) DeleteInstallation(installationId string) error {
	// Delete the installation record
	return r.db.Where("id = ?", installationId).Delete(&entities.CustomerInstallation{}).Error
}
