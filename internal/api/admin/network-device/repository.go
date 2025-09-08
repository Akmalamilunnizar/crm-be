package networkdevice

import (
	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminNetworkDeviceRepositoryInterface interface {
	GetAllAdminNetworkDeviceRepository() ([]entities.NetworkDevice, error)
	GetByIdAdminNetworkDeviceRepository(request IdNetworkDeviceRequest) (entities.NetworkDevice, error)
	GetByCustomerIdAdminNetworkDeviceRepository(customerId string) ([]entities.NetworkDevice, error)
	CreateAdminNetworkDeviceRepository(request CreateNetworkDeviceRequest) (entities.NetworkDevice, error)
	UpdateAdminNetworkDeviceRepository(request UpdateNetworkDeviceRequest) (entities.NetworkDevice, error)
	DeleteAdminNetworkDeviceRepository(request IdNetworkDeviceRequest) (entities.NetworkDevice, error)
}

type AdminNetworkDeviceRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminNetworkDeviceRepository(db *gorm.DB) AdminNetworkDeviceRepositoryInterface {
	return &AdminNetworkDeviceRepositoryStruct{db: db}
}

func (r AdminNetworkDeviceRepositoryStruct) GetAllAdminNetworkDeviceRepository() ([]entities.NetworkDevice, error) {
	var networkDevices []entities.NetworkDevice
	err := r.db.Find(&networkDevices).Error
	return networkDevices, err
}

func (r AdminNetworkDeviceRepositoryStruct) GetByIdAdminNetworkDeviceRepository(request IdNetworkDeviceRequest) (entities.NetworkDevice, error) {
	var networkDevice entities.NetworkDevice
	err := r.db.Where("id = ?", request.Id).First(&networkDevice).Error
	return networkDevice, err
}

func (r AdminNetworkDeviceRepositoryStruct) GetByCustomerIdAdminNetworkDeviceRepository(customerId string) ([]entities.NetworkDevice, error) {
	var networkDevices []entities.NetworkDevice
	err := r.db.Where("customer_id = ?", customerId).Find(&networkDevices).Error
	return networkDevices, err
}

func (r AdminNetworkDeviceRepositoryStruct) CreateAdminNetworkDeviceRepository(request CreateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	assetsID := ""
	if request.AssetsID != nil {
		assetsID = *request.AssetsID
	}

	networkDevice := entities.NetworkDevice{
		ID:              uuid.New().String(),
		CustomerID:      request.CustomerID,
		IPStatic:        request.IPStatic,
		MacAddress:      request.MacAddress,
		StatusPerangkat: request.StatusPerangkat,
		LastPingStatus:  request.LastPingStatus,
		AssetsID:        assetsID,
	}

	err := r.db.Create(&networkDevice).Error
	return networkDevice, err
}

func (r AdminNetworkDeviceRepositoryStruct) UpdateAdminNetworkDeviceRepository(request UpdateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	var networkDevice entities.NetworkDevice
	err := r.db.Where("id = ?", request.ID).First(&networkDevice).Error
	if err != nil {
		return networkDevice, err
	}

	assetsID := ""
	if request.AssetsID != nil {
		assetsID = *request.AssetsID
	}

	networkDevice.CustomerID = request.CustomerID
	networkDevice.IPStatic = request.IPStatic
	networkDevice.MacAddress = request.MacAddress
	networkDevice.StatusPerangkat = request.StatusPerangkat
	networkDevice.LastPingStatus = request.LastPingStatus
	networkDevice.AssetsID = assetsID

	err = r.db.Save(&networkDevice).Error
	return networkDevice, err
}

func (r AdminNetworkDeviceRepositoryStruct) DeleteAdminNetworkDeviceRepository(request IdNetworkDeviceRequest) (entities.NetworkDevice, error) {
	var networkDevice entities.NetworkDevice
	err := r.db.Where("id = ?", request.Id).First(&networkDevice).Error
	if err != nil {
		return networkDevice, err
	}

	err = r.db.Delete(&networkDevice).Error
	return networkDevice, err
}
