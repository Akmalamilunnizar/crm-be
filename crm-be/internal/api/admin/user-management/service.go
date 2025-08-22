package usermanagement

import (
	"skripsi-be/internal/models/dto"
)

type AdminUserManagementServiceInterface interface {
	GetAllAdminUserManagementService(request SearchAdminUserManagementRequest) (*[]dto.UserDTO, error)
	GetByIdAdminUserManagementService(request IdAdminUserManagementRequest) (dto.UserDTO, error)
	CreateAdminUserManagementService(request CreateAdminUserManagementRequest) (dto.UserDTO, error)
	UpdateAdminUserManagementService(request UpdateAdminUserManagementRequest) (dto.UserDTO, error)
	DeleteAdminUserManagementService(request IdAdminUserManagementRequest) (dto.UserDTO, error)
}

type AdminUserManagementServiceStruct struct {
	repository AdminUserManagementRepositoryInterface
}

func NewAdminUserManagementService(repository *AdminUserManagementRepositoryStruct) *AdminUserManagementServiceStruct {
	return &AdminUserManagementServiceStruct{repository}
}

func (s *AdminUserManagementServiceStruct) GetAllAdminUserManagementService(request SearchAdminUserManagementRequest) (*[]dto.UserDTO, error) {
	Users, err := s.repository.FindAdminUserManagementRepository(request)
	if err != nil {
		return nil, err
	}

	return Users, nil
}

func (s *AdminUserManagementServiceStruct) GetByIdAdminUserManagementService(request IdAdminUserManagementRequest) (dto.UserDTO, error) {
	user, err := s.repository.FindByIdAdminUserManagementRepository(request)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *AdminUserManagementServiceStruct) CreateAdminUserManagementService(request CreateAdminUserManagementRequest) (dto.UserDTO, error) {
	user, err := s.repository.CreateAdminUserManagementRepository(request)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *AdminUserManagementServiceStruct) UpdateAdminUserManagementService(request UpdateAdminUserManagementRequest) (dto.UserDTO, error) {
	user, err := s.repository.UpdateAdminUserManagementRepository(request)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *AdminUserManagementServiceStruct) DeleteAdminUserManagementService(request IdAdminUserManagementRequest) (dto.UserDTO, error) {
	user, err := s.repository.DeleteAdminUserManagementRepository(request)
	if err != nil {
		return user, err
	}
	return user, err
}
