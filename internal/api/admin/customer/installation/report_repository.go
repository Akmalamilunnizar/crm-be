package customerinstallation

import (
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type AdminInstallationReportRepositoryInterface interface {
	FindCompleteInstallationReportRepository(installationId string) (entities.CustomerInstallation, error)
	FindAllCompleteInstallationReportsRepository() ([]InstallationReportCompleteResponse, error)
	FindInstallationSummaryPerCustomerRepository() ([]InstallationSummaryResponse, error)
	FindInstallationAssetReportRepository(installationId string) (InstallationAssetReportResponse, error)
	FindInstallationTechnicianReportRepository() ([]InstallationTechnicianReportResponse, error)
	CreateCompleteInstallationReportRepository(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
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
			COUNT(ci.id) as total_installations,
			COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
			COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
			COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
			MAX(ci.on_air_date) as latest_on_air_date,
			MAX(ci.installation_completed_at) as latest_completion_date
		FROM customer c
		LEFT JOIN customer_installations ci ON c.id = ci.customer_id
		GROUP BY c.id, c.name, c.address, c.phone
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
		if parsed, err := time.Parse("2006-01-02 15:04:05", request.InstallationCompletedAt); err == nil {
			installationCompletedAt = &parsed
		}
	}

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
			ID:                     "",
			CustomerID:             request.CustomerId,
			AssetsID:               device.AssetsID,
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
