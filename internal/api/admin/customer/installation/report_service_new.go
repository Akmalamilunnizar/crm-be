package customerinstallation

import (
	"skripsi-be/internal/models/entities"
)

type ReportInstallationServiceInterface interface {
	CreateReportInstallationService(request CreateReportInstallationRequest) (entities.CustomerInstallation, error)
}

type ReportInstallationService struct {
	repository ReportInstallationRepositoryInterface
}

func NewReportInstallationService(repository ReportInstallationRepositoryInterface) ReportInstallationServiceInterface {
	return &ReportInstallationService{repository: repository}
}

// CreateReportInstallationService - Create installation report with all related data
func (s *ReportInstallationService) CreateReportInstallationService(request CreateReportInstallationRequest) (entities.CustomerInstallation, error) {
	installation, err := s.repository.CreateReportInstallationRepository(request)
	if err != nil {
		return installation, err
	}
	return installation, nil
}
