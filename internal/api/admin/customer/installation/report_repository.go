package customerinstallation

import (
	"database/sql"
	"fmt"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type AdminInstallationReportRepositoryInterface interface {
	FindCompleteInstallationReportRepository(installationId string) (entities.CustomerInstallation, error)
	FindCompleteInstallationReportByViewRepository(installationId string) (InstallationReportCompleteResponse, error)
	FindAllCompleteInstallationReportsRepository() ([]InstallationReportCompleteResponse, error)
	FindInstallationSummaryPerCustomerRepository() ([]InstallationSummaryResponse, error)
	FindInstallationAssetReportRepository(installationId string) (InstallationAssetReportResponse, error)
	FindInstallationTechnicianReportRepository() ([]InstallationTechnicianReportResponse, error)
	CreateCompleteInstallationReportRepository(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
	UpdateCompleteInstallationReportRepository(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
}

type AdminInstallationReportRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminInstallationReportRepository(db *gorm.DB) AdminInstallationReportRepositoryStruct {
	return AdminInstallationReportRepositoryStruct{db}
}


// FindCompleteInstallationReportRepository - Get complete installation report with all related data
func (r AdminInstallationReportRepositoryStruct) FindCompleteInstallationReportRepository(installationId string) (entities.CustomerInstallation, error) {
	var installation entities.CustomerInstallation

	err := r.db.Preload("Customer").
		Preload("Technician").
		Preload("Images").
		Preload("AssetTransactions").
		Preload("NetworkDevices").
		Preload("NetworkDevices.Assets").
		Preload("NetworkDevices.Product").
		Preload("CustomerServices").
		Preload("Cables").
		Where("id = ?", installationId).
		First(&installation).Error

	return installation, err
}

// FindCompleteInstallationReportByViewRepository - Get complete installation report using database view
func (r AdminInstallationReportRepositoryStruct) FindCompleteInstallationReportByViewRepository(installationId string) (InstallationReportCompleteResponse, error) {
	var report InstallationReportCompleteResponse

	query := `
		SELECT 
			installation_id,
			customer_id,
			customer_name,
			customer_address,
			customer_phone,
			tgl_permintaan_psb,
			technician_id,
			technician_name,
			technician_phone,
			installation_status,
			installation_type,
			installation_notes,
			on_air_date,
			trial_end_date,
			service_ready_date,
			installation_completed_at,
			durasi_psb,
			status_psb,
			document_type,
			document_photo,
			total_assets_out,
			total_assets_in,
			network_device_id,
			switch_id,
			port_number,
			remote_port,
			eth_port,
			mac_address,
			ip_static,
			status_perangkat,
			kepemilikan_perangkat,
			last_ping_status,
			last_ping_timestamp,
			router_brand,
			router_type,
			router_model,
			router_serial,
			customer_service_id,
			user_login,
			password,
			user_status,
			service_notes,
			cable_id,
			cable_name,
			cable_type,
			cable_length,
			cable_status,
			end_port_type,
			installation_created_at,
			installation_updated_at
		FROM installation_report_complete
		WHERE installation_id = ?
		LIMIT 1
	`

	err := r.db.Raw(query, installationId).Scan(&report).Error
	return report, err
}

// FindAllCompleteInstallationReportsRepository - Get all complete installation reports using database view
func (r AdminInstallationReportRepositoryStruct) FindAllCompleteInstallationReportsRepository() ([]InstallationReportCompleteResponse, error) {
	var reports []InstallationReportCompleteResponse

	// Use the installation_report_complete view to get all reports
	query := `
		SELECT 
			installation_id,
			customer_id,
			customer_name,
			customer_address,
			customer_phone,
			tgl_permintaan_psb,
			technician_id,
			technician_name,
			technician_phone,
			installation_status,
			installation_type,
			installation_notes,
			on_air_date,
			trial_end_date,
			service_ready_date,
			installation_completed_at,
			durasi_psb,
			status_psb,
			document_type,
			document_photo,
			total_assets_out,
			total_assets_in,
			network_device_id,
			switch_id,
			port_number,
			remote_port,
			eth_port,
			mac_address,
			ip_static,
			status_perangkat,
			kepemilikan_perangkat,
			last_ping_status,
			last_ping_timestamp,
			router_brand,
			router_type,
			router_model,
			router_serial,
			customer_service_id,
			user_login,
			password,
			user_status,
			service_notes,
			cable_id,
			cable_name,
			cable_type,
			cable_length,
			cable_status,
			end_port_type,
			installation_created_at,
			installation_updated_at
		FROM installation_report_complete
		ORDER BY installation_created_at DESC
	`

	err := r.db.Raw(query).Scan(&reports).Error
	return reports, err
}

// FindInstallationSummaryPerCustomerRepository - Get installation summary grouped by customer
func (r AdminInstallationReportRepositoryStruct) FindInstallationSummaryPerCustomerRepository() ([]InstallationSummaryResponse, error) {
	var summaries []InstallationSummaryResponse

	query := `
		SELECT 
			c.id as customer_id,
			c.name as customer_name,
			c.address as customer_address,
			c.phone as customer_phone,
			c.service_request_date as tgl_permintaan_psb,
			COUNT(ci.id) as total_installations,
			COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
			COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
			COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
			MAX(ci.on_air_date) as latest_on_air_date,
			MAX(ci.installation_completed_at) as latest_completion_date,
			AVG(CASE 
				WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
				THEN DATEDIFF(ci.installation_completed_at, c.service_request_date)
				ELSE NULL
			END) as avg_durasi_psb,
			COUNT(CASE 
				WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
				AND DATEDIFF(ci.installation_completed_at, c.service_request_date) <= 3
				THEN 1
			END) as tepat_waktu_count,
			COUNT(CASE 
				WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
				AND DATEDIFF(ci.installation_completed_at, c.service_request_date) > 3
				THEN 1
			END) as terlambat_count
		FROM customer c
		LEFT JOIN customer_installations ci ON c.id = ci.customer_id
		GROUP BY c.id, c.name, c.address, c.phone, c.service_request_date
		ORDER BY total_installations DESC
	`

	err := r.db.Raw(query).Scan(&summaries).Error
	return summaries, err
}

// FindInstallationAssetReportRepository - Get asset report for specific installation
func (r AdminInstallationReportRepositoryStruct) FindInstallationAssetReportRepository(installationId string) (InstallationAssetReportResponse, error) {
	var report InstallationAssetReportResponse

	query := `
		SELECT 
			ci.id as installation_id,
			c.name as customer_name,
			ci.installation_type,
			ci.status as installation_status,
			ci.on_air_date,
			ci.installation_completed_at,
			COUNT(CASE WHEN at.transaction_type = 'out' THEN 1 END) as total_assets_out,
			COALESCE(SUM(CASE WHEN at.transaction_type = 'out' THEN at.quantity ELSE 0 END), 0) as total_quantity_out,
			COUNT(CASE WHEN at.transaction_type = 'in' THEN 1 END) as total_assets_in,
			COALESCE(SUM(CASE WHEN at.transaction_type = 'in' THEN at.quantity ELSE 0 END), 0) as total_quantity_in,
			GROUP_CONCAT(DISTINCT CASE WHEN at.transaction_type = 'out' THEN CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)') END SEPARATOR ', ') as assets_out_details,
			GROUP_CONCAT(DISTINCT CASE WHEN at.transaction_type = 'in' THEN CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)') END SEPARATOR ', ') as assets_in_details
		FROM customer_installations ci
		LEFT JOIN customer c ON ci.customer_id = c.id
		LEFT JOIN asset_transactions at ON ci.id = at.customer_installation_id
		LEFT JOIN assets a ON at.asset_id = a.id
		WHERE ci.id = ?
		GROUP BY ci.id, c.name, ci.installation_type, ci.status, ci.on_air_date, ci.installation_completed_at
	`

	err := r.db.Raw(query, installationId).Scan(&report).Error
	return report, err
}

// FindInstallationTechnicianReportRepository - Get technician performance report
func (r AdminInstallationReportRepositoryStruct) FindInstallationTechnicianReportRepository() ([]InstallationTechnicianReportResponse, error) {
	var reports []InstallationTechnicianReportResponse

	query := `
		SELECT 
			u.id as technician_id,
			u.name as technician_name,
			u.phone as technician_phone,
			COUNT(ci.id) as total_installations,
			COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
			COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
			COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
			AVG(TO_DAYS(ci.installation_completed_at) - TO_DAYS(ci.createdAt)) as avg_completion_days,
			MAX(ci.installation_completed_at) as latest_completion_date
		FROM users u
		LEFT JOIN customer_installations ci ON u.id = ci.technician_id
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'TECHNICIAN')
		GROUP BY u.id, u.name, u.phone
		ORDER BY total_installations DESC
	`

	err := r.db.Raw(query).Scan(&reports).Error
	return reports, err
}

// CreateCompleteInstallationReportRepository - Create complete installation report with all related data
func (r AdminInstallationReportRepositoryStruct) CreateCompleteInstallationReportRepository(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
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
			fmt.Printf("Successfully parsed installation_completed_at: %s -> %v\n", request.InstallationCompletedAt, *installationCompletedAt)
		} else {
			fmt.Printf("Failed to parse installation_completed_at: %s\n", request.InstallationCompletedAt)
		}
	}

	// Log document photo for debugging
	fmt.Printf("Document Photo from request: '%s'\n", request.DocumentPhoto)
	fmt.Printf("Document Type from request: '%s'\n", request.DocumentType)

	// Create main installation record
	installation := entities.CustomerInstallation{
		ID:                      "",
		CustomerID:              &request.CustomerId,
		TechnicianID:            &request.TechnicianId,
		Status:                  request.Status,
		Notes:                   request.Notes,
		DocumentType:            &request.DocumentType,
		DocumentPhoto:           &request.DocumentPhoto,
		InstallationType:        request.InstallationType,
		TotalAssetsOut:          request.TotalAssetsOut,
		TotalAssetsIn:           request.TotalAssetsIn,
		InstallationCompletedAt: installationCompletedAt,
		TrialEndDate:            trialEndDate,
		ServiceReadyDate:        serviceReadyDate,
		OnAirDate:               onAirDate,
	}

	if err := tx.Create(&installation).Error; err != nil {
		tx.Rollback()
		return installation, err
	}

	// Log what was actually saved to database
	fmt.Printf("Installation created with ID: %s\n", installation.ID)
	if installation.DocumentPhoto != nil {
		// Normalize the document photo path if needed
		normalizedPath := normalizeDocumentPhotoPath(*installation.DocumentPhoto)
		if normalizedPath != *installation.DocumentPhoto {
			// Update the database with the normalized path
			installation.DocumentPhoto = &normalizedPath
			tx.Save(&installation)
			fmt.Printf("Document Photo normalized from '%s' to '%s'\n", *installation.DocumentPhoto, normalizedPath)
		} else {
			fmt.Printf("Document Photo saved to DB: '%s'\n", *installation.DocumentPhoto)
		}
	} else {
		fmt.Printf("Document Photo is NULL in DB\n")
	}
	if installation.DocumentType != nil {
		fmt.Printf("Document Type saved to DB: '%s'\n", *installation.DocumentType)
	} else {
		fmt.Printf("Document Type is NULL in DB\n")
	}

	// Create asset transactions
	for _, assetTx := range request.AssetTransactions {
		var transactionDate time.Time
		if assetTx.TransactionDate != "" {
			if parsed, err := time.Parse("2006-01-02 15:04:05", assetTx.TransactionDate); err == nil {
				transactionDate = parsed
			} else {
				transactionDate = time.Now()
			}
		} else {
			transactionDate = time.Now()
		}

		assetTransaction := entities.AssetTransaction{
			ID:                     "",
			CustomerInstallationID: installation.ID,
			AssetID:                assetTx.AssetId,
			TransactionType:        assetTx.TransactionType,
			Quantity:               assetTx.Quantity,
			Notes:                  assetTx.Notes,
			TransactionDate:        transactionDate,
			CreatedBy:              request.TechnicianId, // Assuming technician creates the transaction
		}

		if err := tx.Create(&assetTransaction).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create network devices
	for _, device := range request.NetworkDevices {
		networkDevice := entities.NetworkDevice{
			ID:         "",
			CustomerID: request.CustomerId,
			AssetsID: sql.NullString{
				String: func() string {
					if device.AssetsID != "" {
						return device.AssetsID
					}
					return ""
				}(),
				Valid: device.AssetsID != "",
			},
			CustomerInstallationID: &installation.ID,
			SwitchID:               &device.SwitchID,
			PortNumber:             &device.PortNumber,
			RemotePort:             &device.RemotePort,
			EthPort:                &device.EthPort,
			MacAddress:             &device.MacAddress,
			IPStatic:               &device.IPStatic,
			KepemilikanPerangkat:   device.KepemilikanPerangkat,
			StatusPerangkat:        device.StatusPerangkat,
			LastPingStatus:         device.LastPingStatus,
			ProductID:              &device.ProductID,
		}

		if err := tx.Create(&networkDevice).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create customer services
	for _, service := range request.CustomerServices {
		var serviceActivationDate *time.Time
		if service.ServiceActivationDate != "" {
			if parsed, err := time.Parse("2006-01-02", service.ServiceActivationDate); err == nil {
				serviceActivationDate = &parsed
			}
		}

		customerService := entities.CustomerService{
			ID:                     "",
			CustomerID:             request.CustomerId,
			CustomerInstallationID: &installation.ID,
			DeviceID:               &service.DeviceID,
			CableID:                &service.CableID,
			CableLength:            &service.CableLength,
			EndPortType:            &service.EndPortType,
			UserLogin:              &service.UserLogin,
			Password:               &service.Password,
			UserStatus:             service.UserStatus,
			InstallationNotes:      &service.InstallationNotes,
			ServiceActivationDate:  serviceActivationDate,
		}

		if err := tx.Create(&customerService).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create cables
	for _, cable := range request.Cables {
		cableEntity := entities.Cable{
			ID:                     "",
			Name:                   cable.Name,
			Type:                   &cable.Type,
			Length:                 &cable.Length,
			Status:                 cable.Status,
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&cableEntity).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Update images with installation ID
	if len(request.ImageIds) > 0 {
		if err := tx.Model(&entities.Image{}).
			Where("id IN ?", request.ImageIds).
			Update("archive_installation_id", installation.ID).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return installation, err
	}

	return installation, nil
}

// UpdateCompleteInstallationReportRepository - Update complete installation report with all related data
func (r AdminInstallationReportRepositoryStruct) UpdateCompleteInstallationReportRepository(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
	var installation entities.CustomerInstallation

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find existing installation
	if err := tx.Preload("Customer").Preload("Technician").Preload("Images").First(&installation, "id = ?", installationId).Error; err != nil {
		tx.Rollback()
		return installation, err
	}

	// Update basic installation data
	installation.CustomerID = &request.CustomerId
	installation.TechnicianID = &request.TechnicianId
	installation.Status = request.Status
	installation.Notes = request.Notes
	if request.DocumentType != "" {
		installation.DocumentType = &request.DocumentType
	}
	if request.DocumentPhoto != "" {
		installation.DocumentPhoto = &request.DocumentPhoto
	}
	installation.InstallationType = request.InstallationType

	// Parse dates
	if request.OnAirDate != "" {
		if onAirDate, err := time.Parse("2006-01-02", request.OnAirDate); err == nil {
			installation.OnAirDate = &onAirDate
		}
	}
	if request.TrialEndDate != "" {
		if trialEndDate, err := time.Parse("2006-01-02", request.TrialEndDate); err == nil {
			installation.TrialEndDate = &trialEndDate
		}
	}
	if request.ServiceReadyDate != "" {
		if serviceReadyDate, err := time.Parse("2006-01-02", request.ServiceReadyDate); err == nil {
			installation.ServiceReadyDate = &serviceReadyDate
		}
	}
	if request.InstallationCompletedAt != "" {
		if completedAt, err := time.Parse("2006-01-02T15:04:05", request.InstallationCompletedAt); err == nil {
			installation.InstallationCompletedAt = &completedAt
		}
	}

	// Save installation
	if err := tx.Save(&installation).Error; err != nil {
		tx.Rollback()
		return installation, err
	}

	// Delete existing related data
	if err := tx.Where("customer_installation_id = ?", installation.ID).Delete(&entities.NetworkDevice{}).Error; err != nil {
		tx.Rollback()
		return installation, err
	}
	if err := tx.Where("customer_installation_id = ?", installation.ID).Delete(&entities.CustomerService{}).Error; err != nil {
		tx.Rollback()
		return installation, err
	}
	if err := tx.Where("customer_installation_id = ?", installation.ID).Delete(&entities.Cable{}).Error; err != nil {
		tx.Rollback()
		return installation, err
	}

	// Create network devices
	for _, device := range request.NetworkDevices {
		deviceEntity := entities.NetworkDevice{
			ID:                     "",
			CustomerID:             request.CustomerId,
			AssetsID:               &device.AssetsID,
			SwitchID:               &device.SwitchID,
			PortNumber:             &device.PortNumber,
			RemotePort:             &device.RemotePort,
			EthPort:                &device.EthPort,
			MacAddress:             &device.MacAddress,
			IPStatic:               &device.IPStatic,
			KepemilikanPerangkat:   device.KepemilikanPerangkat,
			StatusPerangkat:        device.StatusPerangkat,
			LastPingStatus:         device.LastPingStatus,
			ProductID:              &device.ProductID,
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&deviceEntity).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create customer services
	for _, service := range request.CustomerServices {
		var serviceActivationDate *time.Time
		if service.ServiceActivationDate != "" {
			if parsedDate, err := time.Parse("2006-01-02", service.ServiceActivationDate); err == nil {
				serviceActivationDate = &parsedDate
			}
		}

		serviceEntity := entities.CustomerService{
			ID:                     "",
			DeviceID:               &service.DeviceID,
			CableID:                &service.CableID,
			CableLength:            &service.CableLength,
			EndPortType:            &service.EndPortType,
			UserLogin:              &service.UserLogin,
			Password:               &service.Password,
			UserStatus:             service.UserStatus,
			InstallationNotes:      &service.InstallationNotes,
			ServiceActivationDate:  serviceActivationDate,
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&serviceEntity).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Create cables
	for _, cable := range request.Cables {
		cableEntity := entities.Cable{
			ID:                     "",
			Name:                   cable.Name,
			Type:                   &cable.Type,
			Length:                 &cable.Length,
			Status:                 cable.Status,
			CustomerInstallationID: &installation.ID,
		}

		if err := tx.Create(&cableEntity).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	// Update images with installation ID
	if len(request.ImageIds) > 0 {
		if err := tx.Model(&entities.Image{}).
			Where("id IN ?", request.ImageIds).
			Update("archive_installation_id", installation.ID).Error; err != nil {
			tx.Rollback()
			return installation, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return installation, err
	}

	return installation, nil
}
