package company

import "skripsi-be/internal/models/entities"

type AdminCompanyServiceInterface interface {
}

type AdminCompanyServiceStruct struct {
	repository *AdminCompanyRepositoryStruct
}

func NewAdminCompanyService(repository *AdminCompanyRepositoryStruct) *AdminCompanyServiceStruct {
	return &AdminCompanyServiceStruct{repository}
}

func (s *AdminCompanyServiceStruct) GetAllAdminCompanyService() (*[]entities.Company, error) {
	companies, err := s.repository.GetAllAdminCompanyRepository()
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (s *AdminCompanyServiceStruct) GetByIdAdminCompanyService(request IdAdminCompanyRequest) (*entities.Company, error) {
	companies, err := s.repository.GetByIDAdminCompanyRepository(request)
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (s *AdminCompanyServiceStruct) CreateAdminCompanyService(request CreateAdminCompanyRequest) {
	s.repository.CreateAdminCompanyRepository(request)
}

func (s *AdminCompanyServiceStruct) UpdateAdminCompanyService(request UpdateAdminCompanyRequest) {
	s.repository.UpdateAdminCompanyRepository(request)
}

func (s *AdminCompanyServiceStruct) DeleteAdminCompanyService(request IdAdminCompanyRequest) {
	s.repository.DeleteCompanyByID(request)
}
