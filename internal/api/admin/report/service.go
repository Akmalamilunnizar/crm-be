package report

import "skripsi-be/internal/models/dto"

type AdminReportServiceInterface interface {
	GetAllInternetAdminReportService() ([]dto.ReportInternet, error)
}
type AdminReportServiceStruct struct {
	repository AdminReportRepositoryInterface
}

func NewAdminReportService(repository AdminReportRepositoryInterface) AdminReportServiceStruct {
	return AdminReportServiceStruct{repository}
}

func (s AdminReportServiceStruct) GetAllInternetAdminReportService() ([]dto.ReportInternet, error) {
	customers, err := s.repository.InternetReportAdminReportRepository()
	if err != nil {
		return nil, err
	}

	return customers, nil
}
