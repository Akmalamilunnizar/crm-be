package account

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
)

type AdminAccountServiceInterface interface {
	GetAllAdminAccountService() ([]entities.Accounts, error)
	CreateAdminAccountService(request CreateAdminAccountRequest) (entities.Accounts, error)
	GetByIdAdminAccountService(request IdAdminAccountRequest) (entities.Accounts, error)
	UpdateAdminAccountService(request UpdateAdminAccountRequest) (entities.Accounts, error)
	DeleteAdminAccountService(request IdAdminAccountRequest) (entities.Accounts, error)
}
type AdminAccountServiceStruct struct {
	repository AdminAccountRepositoryInterface
}

func NewAdminAccountService(repository AdminAccountRepositoryInterface) *AdminAccountServiceStruct {
	return &AdminAccountServiceStruct{repository}
}

func (s AdminAccountServiceStruct) GetAllAdminAccountService() ([]entities.Accounts, error) {
	areas, err := s.repository.FindAdminAccountRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s AdminAccountServiceStruct) GetByIdAdminAccountService(request IdAdminAccountRequest) (entities.Accounts, error) {
	area := entities.Accounts{}
	err := copier.Copy(&area, &request)
	if err != nil {
		return area, err
	}

	area, err = s.repository.FindByIdAdminAccountRepository(request)
	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAccountServiceStruct) CreateAdminAccountService(request CreateAdminAccountRequest) (entities.Accounts, error) {

	area, err := s.repository.CreateAdminAccountRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAccountServiceStruct) UpdateAdminAccountService(request UpdateAdminAccountRequest) (entities.Accounts, error) {
	area, err := s.repository.UpdateAdminAccountRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAccountServiceStruct) DeleteAdminAccountService(request IdAdminAccountRequest) (entities.Accounts, error) {
	area, err := s.repository.DeleteAdminAccountRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}
