package customer

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/copier"

	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
)

type AdminCustomerServiceInterface interface {
	GetAllAdminCustomerService() ([]entities.Customer, error)
	GetByIdAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error)
	CreateAdminCustomerService(request CreateAdminCustomerRequest) (entities.Customer, error)
	UpdateAdminCustomerService(request UpdateAdminCustomerRequest) (entities.Customer, error)
	DeleteAdminCustomerService(request IdAdminCustomerRequest) (entities.Customer, error)
}

type AdminCustomerServiceStruct struct {
	repository AdminCustomerRepositoryInterface
}

func NewAdminCustomerService(repository AdminCustomerRepositoryInterface) AdminCustomerServiceStruct {
	return AdminCustomerServiceStruct{repository}
}

func (s AdminCustomerServiceStruct) GetAllAdminCustomerService() ([]entities.Customer, error) {
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
