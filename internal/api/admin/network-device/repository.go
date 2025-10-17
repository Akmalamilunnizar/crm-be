package networkdevice

import (
	"database/sql"
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
	err := r.db.Preload("Product").Preload("Customer", "deleted_at IS NULL").Find(&networkDevices).Error
	return networkDevices, err
}

func (r AdminNetworkDeviceRepositoryStruct) GetByIdAdminNetworkDeviceRepository(request IdNetworkDeviceRequest) (entities.NetworkDevice, error) {
	var networkDevice entities.NetworkDevice
	err := r.db.Preload("Product").Preload("Customer", "deleted_at IS NULL").Where("id = ?", request.Id).First(&networkDevice).Error
	return networkDevice, err
}

func (r AdminNetworkDeviceRepositoryStruct) GetByCustomerIdAdminNetworkDeviceRepository(customerId string) ([]entities.NetworkDevice, error) {
	var networkDevices []entities.NetworkDevice
	err := r.db.Preload("Product").Preload("Customer", "deleted_at IS NULL").Where("customer_id = ?", customerId).Find(&networkDevices).Error
	return networkDevices, err
}

func (r AdminNetworkDeviceRepositoryStruct) CreateAdminNetworkDeviceRepository(request CreateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	networkDevice := entities.NetworkDevice{
		ID:         uuid.New().String(),
		CustomerID: request.CustomerID,
		IPStatic:   request.IPStatic,
		MacAddress: request.MacAddress,
		AssetsID: sql.NullString{
			String: func() string {
				if request.AssetsID != nil {
					return *request.AssetsID
				}
				return ""
			}(),
			Valid: request.AssetsID != nil && *request.AssetsID != "",
		},
		ProductID: request.ProductID,
	}

	err := r.db.Create(&networkDevice).Error
	if err != nil {
		return networkDevice, err
	}

	// Preload the product information after creation
	err = r.db.Preload("Product").Preload("Customer").Where("id = ?", networkDevice.ID).First(&networkDevice).Error
	return networkDevice, err
}

func (r AdminNetworkDeviceRepositoryStruct) UpdateAdminNetworkDeviceRepository(request UpdateNetworkDeviceRequest) (entities.NetworkDevice, error) {
	var networkDevice entities.NetworkDevice
	err := r.db.Where("id = ?", request.ID).First(&networkDevice).Error
	if err != nil {
		return networkDevice, err
	}

	var assetsID *string
	if request.AssetsID != nil && *request.AssetsID != "" {
		assetsID = request.AssetsID
	}

	networkDevice.CustomerID = request.CustomerID
	networkDevice.IPStatic = request.IPStatic
	networkDevice.MacAddress = request.MacAddress
	networkDevice.AssetsID = sql.NullString{
		String: func() string {
			if assetsID != nil {
				return *assetsID
			}
			return ""
		}(),
		Valid: assetsID != nil && *assetsID != "",
	}
	networkDevice.ProductID = request.ProductID

	err = r.db.Save(&networkDevice).Error
	if err != nil {
		return networkDevice, err
	}

	// Preload the product information after update
	err = r.db.Preload("Product").Preload("Customer").Where("id = ?", networkDevice.ID).First(&networkDevice).Error
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
