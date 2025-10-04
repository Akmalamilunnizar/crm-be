package customerinstallation

import (
	"skripsi-be/internal/models/entities"
)

type AdminInstallationReportServiceInterface interface {
	GetCompleteInstallationReportService(installationId string) (entities.CustomerInstallation, error)
	GetCompleteInstallationReportByViewService(installationId string) (InstallationReportCompleteResponse, error)
	GetAllCompleteInstallationReportsService() ([]InstallationReportCompleteResponse, error)
	GetInstallationSummaryPerCustomerService() ([]InstallationSummaryResponse, error)
	GetInstallationAssetReportService(installationId string) (InstallationAssetReportResponse, error)
	GetInstallationTechnicianReportService() ([]InstallationTechnicianReportResponse, error)
	CreateCompleteInstallationReportService(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
	UpdateCompleteInstallationReportService(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error)
}

type AdminInstallationReportServiceStruct struct {
	repository AdminInstallationReportRepositoryInterface
}

func NewAdminInstallationReportService(repository AdminInstallationReportRepositoryInterface) AdminInstallationReportServiceStruct {
	return AdminInstallationReportServiceStruct{repository}
}

// GetCompleteInstallationReportService - Get complete installation report with all related data
func (s AdminInstallationReportServiceStruct) GetCompleteInstallationReportService(installationId string) (entities.CustomerInstallation, error) {
	installation, err := s.repository.FindCompleteInstallationReportRepository(installationId)
	if err != nil {
		return installation, err
	}
	return installation, nil
}

// GetCompleteInstallationReportByViewService - Get complete installation report using database view
func (s AdminInstallationReportServiceStruct) GetCompleteInstallationReportByViewService(installationId string) (InstallationReportCompleteResponse, error) {
	report, err := s.repository.FindCompleteInstallationReportByViewRepository(installationId)
	if err != nil {
		return report, err
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

// CreateCompleteInstallationReportService - Create complete installation report with all related data
func (s AdminInstallationReportServiceStruct) CreateCompleteInstallationReportService(request CreateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
	installation, err := s.repository.CreateCompleteInstallationReportRepository(request)
	if err != nil {
		return installation, err
	}
	return installation, nil
}

// UpdateCompleteInstallationReportService - Update complete installation report with all related data
func (s AdminInstallationReportServiceStruct) UpdateCompleteInstallationReportService(installationId string, request UpdateCompleteInstallationReportRequest) (entities.CustomerInstallation, error) {
	installation, err := s.repository.UpdateCompleteInstallationReportRepository(installationId, request)
	if err != nil {
		return installation, err
	}
	return installation, nil
}
