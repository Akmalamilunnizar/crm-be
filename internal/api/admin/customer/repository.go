package customer

import (
	"fmt"
	"strings"
	"time"
	"skripsi-be/internal/models/entities"
	"skripsi-be/internal/services"

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
	tx := r.db.Preload("Area").Find(&customers)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Build response with network device data
	var response []CustomerListResponse
	for _, customer := range customers {
		// Get network devices for this customer (with Product)
		networkDevices := []entities.NetworkDevice{}
		r.db.Preload("Product").Where("customer_id = ?", customer.ID).Find(&networkDevices)

		// Note: Product information is now managed through network_devices table
		// Customer no longer has direct Product relationship

		// Aggregate IP and MAC addresses from network devices
		var ipAddresses []string
		var macAddresses []string
		var productName *string
		var productPrice *int64

		for _, device := range networkDevices {
			if device.IPStatic != nil && *device.IPStatic != "" {
				ipAddresses = append(ipAddresses, *device.IPStatic)
			}
			if device.MacAddress != nil && *device.MacAddress != "" {
				macAddresses = append(macAddresses, *device.MacAddress)
			}
			// Get product information from the first device with a product
			if device.Product != nil && productName == nil {
				productName = &device.Product.Name
				productPrice = &device.Product.Price
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
			Customer:     customer,
			IPStatic:     combinedIP,
			MacAddress:   combinedMAC,
			ProductName:  productName,
			ProductPrice: productPrice,
		}

		response = append(response, customerResponse)
	}

	return response, nil
}
func (r AdminCustomerRepositoryStruct) FindByIdAdminCustomerRepository(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer := entities.Customer{}
	tx := r.db.Preload("Area").Find(&customer, "id = ?", request.Id)
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
	tx := r.db.Preload("Area").Find(&customer, "id = ?", request.Id)
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
	tx := r.db.Preload("Area").Find(&customer, "id = ?", request.Id)
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
	IPStatic     *string `json:"ip_static"`
	MacAddress   *string `json:"mac_address"`
	ProductName  *string `json:"product_name"`
	ProductPrice *int64  `json:"product_price"`
}

func (r AdminCustomerRepositoryStruct) FindByIdDetailAdminCustomerRepository(request IdAdminCustomerRequest) (*CustomerDetailResponse, error) {
	// Get customer with all related data
	customer := entities.Customer{}
	tx := r.db.Preload("Area").First(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Get customer installations
	installations := []entities.CustomerInstallation{}
	r.db.Preload("Technician").Preload("Images").Where("customer_id = ?", request.Id).Find(&installations)

	// Get network devices for this customer (with Product)
	networkDevices := []entities.NetworkDevice{}
	r.db.Preload("Product").Where("customer_id = ?", request.Id).Find(&networkDevices)

	// Get real-time Netwatch status for each network device from MikroTik
	for i := range networkDevices {
		if networkDevices[i].IPStatic != nil && *networkDevices[i].IPStatic != "" {
			// Get real-time status from MikroTik Netwatch
			status, lastSeen, err := r.getRealTimeNetwatchStatus(*networkDevices[i].IPStatic)
			if err == nil {
				// Update the network device's ping status with real-time MikroTik data
				networkDevices[i].LastPingStatus = status
				networkDevices[i].LastPingTimestamp = &lastSeen
			} else {
				// Fallback to database if MikroTik is not available
				var netwatchDevice entities.NetwatchDevice
				tx := r.db.Where("ip_address = ?", *networkDevices[i].IPStatic).First(&netwatchDevice)
				if tx.Error == nil {
					networkDevices[i].LastPingStatus = netwatchDevice.Status
					networkDevices[i].LastPingTimestamp = &netwatchDevice.LastSeen
				}
			}
		}
	}

	// Note: Product information is now managed through network_devices table
	// Customer no longer has direct Product relationship

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

// getRealTimeNetwatchStatus gets real-time status from MikroTik Netwatch
func (r AdminCustomerRepositoryStruct) getRealTimeNetwatchStatus(ipAddress string) (string, time.Time, error) {
	// Use shared MikroTik service that's already connected
	mikroTikService := services.GetSharedMikroTikService()
	if mikroTikService == nil || !mikroTikService.IsConnected() {
		return "", time.Time{}, fmt.Errorf("MikroTik service not available or not connected")
	}

	// Get Netwatch devices from MikroTik
	devices, err := mikroTikService.GetNetwatchDevices()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get netwatch devices: %v", err)
	}

	// Find the specific device by IP
	for _, device := range devices {
		if host, ok := device["host"].(string); ok && host == ipAddress {
			if status, ok := device["status"].(string); ok {
				// Convert status to lowercase for consistency
				status = strings.ToLower(status)
				return status, time.Now(), nil
			}
		}
	}

	return "unknown", time.Time{}, fmt.Errorf("device %s not found in MikroTik Netwatch", ipAddress)
}
