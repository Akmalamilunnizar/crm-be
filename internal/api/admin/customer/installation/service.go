package customerinstallation

import (
	"skripsi-be/internal/models/entities"
)

type AdminCustomerInstallationServiceInterface interface {
	GetAllAdminCustomerInstallationService() ([]entities.CustomerInstallation, error)
	GetByIdAdminCustomerInstallationService(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	CreateAdminCustomerInstallationService(request CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	UpdateAdminCustomerInstallationService(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	DeleteAdminCustomerInstallationService(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
}

type AdminCustomerInstallationServiceStruct struct {
	repository AdminCustomerInstallationRepositoryInterface
}

func NewAdminCustomerInstallationService(repository AdminCustomerInstallationRepositoryInterface) AdminCustomerInstallationServiceStruct {
	return AdminCustomerInstallationServiceStruct{repository}
}

func (s AdminCustomerInstallationServiceStruct) GetAllAdminCustomerInstallationService() ([]entities.CustomerInstallation, error) {
	customer, err := s.repository.FindAdminCustomerInstallationRepository()
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerInstallationServiceStruct) GetByIdAdminCustomerInstallationService(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer, err := s.repository.FindByIdAdminCustomerInstallationRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerInstallationServiceStruct) CreateAdminCustomerInstallationService(request CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer, err := s.repository.CreateAdminCustomerInstallationRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err

}

func (s AdminCustomerInstallationServiceStruct) UpdateAdminCustomerInstallationService(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer, err := s.repository.UpdateAdminCustomerInstallationRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerInstallationServiceStruct) DeleteAdminCustomerInstallationService(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer, err := s.repository.DeleteAdminCustomerInstallationRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}
