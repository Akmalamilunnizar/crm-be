package customer

import (
	"github.com/jinzhu/copier"

	"skripsi-be/internal/models/entities"
)

type AdminCustomerServiceInterface interface {
	GetAllAdminCustomerService() ([]CustomerListResponse, error)
	GetByIdAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error)
	GetByIdDetailAdminCustomerService(request IdAdminCustomerRequest) (*CustomerDetailResponse, error)
	CreateAdminCustomerService(request CreateAdminCustomerRequest) (entities.Customer, error)
	UpdateAdminCustomerService(request UpdateAdminCustomerRequest) (entities.Customer, error)
	DeleteAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error)
	DeleteAdminCustomerWithRelatedService(request IdAdminCustomerRequest) (entities.Customer, error)
	GetCustomerRelatedRecordsService(request IdAdminCustomerRequest) (map[string]int64, error)
}

type AdminCustomerServiceStruct struct {
	repository AdminCustomerRepositoryInterface
}

func NewAdminCustomerService(repository AdminCustomerRepositoryInterface) AdminCustomerServiceStruct {
	return AdminCustomerServiceStruct{repository}
}

func (s AdminCustomerServiceStruct) GetAllAdminCustomerService() ([]CustomerListResponse, error) {
	customer, err := s.repository.FindAdminCustomerRepository()
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerServiceStruct) GetByIdAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer, err := s.repository.FindByIdAdminCustomerRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerServiceStruct) GetByIdDetailAdminCustomerService(request IdAdminCustomerRequest) (*CustomerDetailResponse, error) {
	detail, err := s.repository.FindByIdDetailAdminCustomerRepository(request)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (s AdminCustomerServiceStruct) CreateAdminCustomerService(request CreateAdminCustomerRequest) (entities.Customer, error) {
	type PayloadMikrotik struct {
		MacAddress string `json:"mac-address"`
		ToAddress  string `json:"to-address"`
		Address    string `json:"address"`
		Type       string `json:"type"`
	}

	customer := entities.Customer{}

	err := copier.Copy(&customer, &request)
	if err != nil {
		return customer, err
	}

	// Note: StatusUser field removed as it doesn't exist in the Customer entity
	// Customer status is managed through customer_services table if needed

	customer, err = s.repository.CreateAdminCustomerRepository(customer)
	if err != nil {
		return customer, err
	}

	// payload := `{
	// "mac-address": "` + customer.MacAddress + `",
	// "to-address": "` + customer.Address + `",
	// "address": "` + customer.Address + `",
	// "type": "bypassed"
	// }`

	// TODO: Update Mikrotik integration to work with network_devices table
	// For now, commenting out to fix database schema mismatch
	/*
		payloadStruct := PayloadMikrotik{
			MacAddress: customer.MacAddress,
			ToAddress:  customer.Address,
			Address:    customer.Address,
			Type:       "bypassed",
		}

		// Convert to JSON
		payload, err := json.Marshal(payloadStruct)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return customer, err
		}

		helpers.HttpRequestHelpers("http://10.3.2.33/rest/ip/hotspot/ip-binding", "PUT", string(payload))
	*/

	return customer, err

}

func (s AdminCustomerServiceStruct) UpdateAdminCustomerService(request UpdateAdminCustomerRequest) (entities.Customer, error) {
	customer, err := s.repository.UpdateAdminCustomerRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerServiceStruct) DeleteAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer, err := s.repository.DeleteAdminCustomerRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerServiceStruct) DeleteAdminCustomerWithRelatedService(request IdAdminCustomerRequest) (entities.Customer, error) {
	customer, err := s.repository.DeleteAdminCustomerWithRelatedRepository(request)
	if err != nil {
		return customer, err
	}

	return customer, err
}

func (s AdminCustomerServiceStruct) GetCustomerRelatedRecordsService(request IdAdminCustomerRequest) (map[string]int64, error) {
	relatedCounts, err := s.repository.GetCustomerRelatedRecordsRepository(request)
	if err != nil {
		return relatedCounts, err
	}

	return relatedCounts, err
}
