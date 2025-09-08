package networkdevice

import (
	"skripsi-be/internal/models/entities"
)

type AdminNetworkDeviceServiceInterface interface {
	GetAllAdminNetworkDeviceService() ([]entities.NetworkDevice, error)
	GetByIdAdminNetworkDeviceService(request IdNetworkDeviceRequest) (entities.NetworkDevice, error)
	GetByCustomerIdAdminNetworkDeviceService(customerId string) ([]entities.NetworkDevice, error)
	CreateAdminNetworkDeviceService(request CreateNetworkDeviceRequest) (entities.NetworkDevice, error)
	UpdateAdminNetworkDeviceService(request UpdateNetworkDeviceRequest) (entities.NetworkDevice, error)
	DeleteAdminNetworkDeviceService(request IdNetworkDeviceRequest) (entities.NetworkDevice, error)
}

type AdminNetworkDeviceServiceStruct struct {
	repository AdminNetworkDeviceRepositoryInterface
}

func NewAdminNetworkDeviceService(repository AdminNetworkDeviceRepositoryInterface) AdminNetworkDeviceServiceInterface {
	return &AdminNetworkDeviceServiceStruct{repository: repository}
}

func (s AdminNetworkDeviceServiceStruct) GetAllAdminNetworkDeviceService() ([]entities.NetworkDevice, error) {
	return s.repository.GetAllAdminNetworkDeviceRepository()
}

func (s AdminNetworkDeviceServiceStruct) GetByIdAdminNetworkDeviceService(request IdNetworkDeviceRequest) (entities.NetworkDevice, error) {
	return s.repository.GetByIdAdminNetworkDeviceRepository(request)
}

func (s AdminNetworkDeviceServiceStruct) GetByCustomerIdAdminNetworkDeviceService(customerId string) ([]entities.NetworkDevice, error) {
	return s.repository.GetByCustomerIdAdminNetworkDeviceRepository(customerId)
}

func (s AdminNetworkDeviceServiceStruct) CreateAdminNetworkDeviceService(request CreateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	return s.repository.CreateAdminNetworkDeviceRepository(request)
}

func (s AdminNetworkDeviceServiceStruct) UpdateAdminNetworkDeviceService(request UpdateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	return s.repository.UpdateAdminNetworkDeviceRepository(request)
}

func (s AdminNetworkDeviceServiceStruct) DeleteAdminNetworkDeviceService(request IdNetworkDeviceRequest) (entities.NetworkDevice, error) {
	return s.repository.DeleteAdminNetworkDeviceRepository(request)
}
