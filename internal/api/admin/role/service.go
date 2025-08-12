package role

import (
	"skripsi-be/internal/models/entities"
)

type AdminRoleServiceInterface interface {
	GetAllAdminRoleService() ([]entities.Role, error)
	CreateAdminRoleService(request CreateAdminRoleRequest) (entities.Role, error)
	GetByIdAdminRoleService(request IdAdminRoleRequest) (entities.Role, error)
	UpdateAdminRoleService(request UpdateAdminRoleRequest) (entities.Role, error)
	DeleteAdminRoleService(request IdAdminRoleRequest) (entities.Role, error)
}
type AdminRoleServiceStruct struct {
	repository AdminRoleRepositoryInterface
}

func NewAdminRoleService(repository AdminRoleRepositoryInterface) AdminRoleServiceStruct {
	return AdminRoleServiceStruct{repository}
}

func (s AdminRoleServiceStruct) GetAllAdminRoleService() ([]entities.Role, error) {
	roles, err := s.repository.FindAdminRoleRepository()
	if err != nil {
		return roles, err
	}

	return roles, nil
}
func (s AdminRoleServiceStruct) GetByIdAdminRoleService(request IdAdminRoleRequest) (entities.Role, error) {
	role, err := s.repository.FindByIdAdminRoleRepository(request)
	if err != nil {
		return role, err
	}
	return role, nil
}
func (s AdminRoleServiceStruct) CreateAdminRoleService(request CreateAdminRoleRequest) (entities.Role, error) {
	role, err := s.repository.CreateAdminRoleRepository(request)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (s AdminRoleServiceStruct) UpdateAdminRoleService(request UpdateAdminRoleRequest) (entities.Role, error) {
	role, err := s.repository.UpdateAdminRoleRepository(request)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (s AdminRoleServiceStruct) DeleteAdminRoleService(request IdAdminRoleRequest) (entities.Role, error) {
	role, err := s.repository.DeleteAdminRoleRepository(request)
	if err != nil {
		return role, err
	}
	return role, nil
}
