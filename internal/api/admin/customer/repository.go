package customer

import (
	"skripsi-be/internal/models/entities"
	"strings"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminCustomerRepositoryInterface interface {
	CreateAdminCustomerRepository(customer entities.Customer) (entities.Customer, error)
	UpdateAdminCustomerRepository(request UpdateAdminCustomerRequest) (entities.Customer, error)
	DeleteAdminCustomerRepository(request IdAdminCustomerRequest) (entities.Customer, error)
	FindByIdAdminCustomerRepository(request IdAdminCustomerRequest) (entities.Customer, error)
	FindByIdDetailAdminCustomerRepository(request IdAdminCustomerRequest) (*CustomerDetailResponse, error)
	FindAdminCustomerRepository() ([]CustomerListResponse, error)
}
type AdminCustomerRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminCustomerRepository(db *gorm.DB) AdminCustomerRepositoryStruct {
	return AdminCustomerRepositoryStruct{db}
}

func (r AdminCustomerRepositoryStruct) FindAdminCustomerRepository() ([]CustomerListResponse, error) {
	customers := []entities.Customer{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customers)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Build response with network device data
	var response []CustomerListResponse
	for _, customer := range customers {
		// Get network devices for this customer (with Product)
		networkDevices := []entities.NetworkDevice{}
		r.db.Preload("Product").Where("customer_id = ?", customer.ID).Find(&networkDevices)

		// If product is managed on network_devices, prefer first device with a loaded product
		for _, device := range networkDevices {
			if device.Product != nil && device.Product.ID != "" {
				customer.Product = device.Product
				customer.ProductID = device.Product.ID
				break
			}
			if device.ProductID != nil && *device.ProductID != "" {
				customer.ProductID = *device.ProductID
				break
			}
		}

		// Aggregate IP and MAC addresses from network devices
		var ipAddresses []string
		var macAddresses []string

		for _, device := range networkDevices {
			if device.IPStatic != nil && *device.IPStatic != "" {
				ipAddresses = append(ipAddresses, *device.IPStatic)
			}
			if device.MacAddress != nil && *device.MacAddress != "" {
				macAddresses = append(macAddresses, *device.MacAddress)
			}
		}

		// Create combined strings
		var combinedIP *string
		var combinedMAC *string

		if len(ipAddresses) > 0 {
			combined := strings.Join(ipAddresses, ", ")
			combinedIP = &combined
		}

		if len(macAddresses) > 0 {
			combined := strings.Join(macAddresses, ", ")
			combinedMAC = &combined
		}

		customerResponse := CustomerListResponse{
			Customer:   customer,
			IPStatic:   combinedIP,
			MacAddress: combinedMAC,
		}

		response = append(response, customerResponse)
	}

	return response, nil
}
func (r AdminCustomerRepositoryStruct) FindByIdAdminCustomerRepository(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer := entities.Customer{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}
func (r AdminCustomerRepositoryStruct) CreateAdminCustomerRepository(customer entities.Customer) (entities.Customer, error) {
	tx := r.db.Create(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, nil

}
func (r AdminCustomerRepositoryStruct) UpdateAdminCustomerRepository(request UpdateAdminCustomerRequest) (entities.Customer, error) {
	customer := entities.Customer{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Save(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}

func (r AdminCustomerRepositoryStruct) DeleteAdminCustomerRepository(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer := entities.Customer{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Delete(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}

// CustomerDetailResponse represents comprehensive customer information
type CustomerDetailResponse struct {
	Customer       entities.Customer               `json:"customer"`
	Installations  []entities.CustomerInstallation `json:"installations"`
	NetworkDevices []entities.NetworkDevice        `json:"network_devices"`
	Invoices       []entities.Invoice              `json:"invoices,omitempty"`
}

// CustomerListResponse represents customer with network device data for list view
type CustomerListResponse struct {
	entities.Customer
	IPStatic   *string `json:"ip_static"`
	MacAddress *string `json:"mac_address"`
}

func (r AdminCustomerRepositoryStruct) FindByIdDetailAdminCustomerRepository(request IdAdminCustomerRequest) (*CustomerDetailResponse, error) {
	// Get customer with all related data
	customer := entities.Customer{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").First(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Get customer installations
	installations := []entities.CustomerInstallation{}
	r.db.Preload("Technician").Preload("Images").Where("customer_id = ?", request.Id).Find(&installations)

	// Get network devices for this customer (with Product)
	networkDevices := []entities.NetworkDevice{}
	r.db.Preload("Product").Where("customer_id = ?", request.Id).Find(&networkDevices)

	// Prefer product from network_devices if present
	for _, device := range networkDevices {
		if device.Product != nil && device.Product.ID != "" {
			customer.Product = device.Product
			customer.ProductID = device.Product.ID
			break
		}
		if device.ProductID != nil && *device.ProductID != "" {
			customer.ProductID = *device.ProductID
			break
		}
	}

	// Get recent invoices for this customer (optional)
	invoices := []entities.Invoice{}
	r.db.Where("customer_id = ?", request.Id).Order("createdAt DESC").Limit(10).Find(&invoices)

	detail := &CustomerDetailResponse{
		Customer:       customer,
		Installations:  installations,
		NetworkDevices: networkDevices,
		Invoices:       invoices,
	}

	return detail, nil
}
