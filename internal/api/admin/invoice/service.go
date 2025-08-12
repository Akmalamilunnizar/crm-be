package invoice

import (
	"github.com/jinzhu/copier"

	"skripsi-be/internal/models/entities"
)

type AdminInvoiceServiceInterface interface {
	GetAllAdminInvoiceService() ([]entities.Invoice, error)
	CreateAdminInvoiceService(request CreateAdminInvoiceRequest) (entities.Invoice, error)
	GetByIdAdminInvoiceService(request IdAdminInvoiceRequest) (entities.Invoice, error)
	UpdateAdminInvoiceService(request UpdateAdminInvoiceRequest) (entities.Invoice, error)
	UpdateStatusAdminInvoiceService(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error)
	DeleteAdminInvoiceService(request IdAdminInvoiceRequest) (entities.Invoice, error)
}
type AdminInvoiceServiceStruct struct {
	repository AdminInvoiceRepositoryInterface
}

func NewAdminInvoiceService(repository AdminInvoiceRepositoryInterface) *AdminInvoiceServiceStruct {
	return &AdminInvoiceServiceStruct{repository}
}

func (s AdminInvoiceServiceStruct) GetAllAdminInvoiceService() ([]entities.Invoice, error) {
	areas, err := s.repository.FindAdminInvoiceRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s AdminInvoiceServiceStruct) GetByIdAdminInvoiceService(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	area := entities.Invoice{}
	err := copier.Copy(&area, &request)
	if err != nil {
		return area, err
	}

	area, err = s.repository.FindByIdAdminInvoiceRepository(request)
	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminInvoiceServiceStruct) CreateAdminInvoiceService(request CreateAdminInvoiceRequest) (entities.Invoice, error) {

	area, err := s.repository.CreateAdminInvoiceRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminInvoiceServiceStruct) UpdateAdminInvoiceService(request UpdateAdminInvoiceRequest) (entities.Invoice, error) {
	area, err := s.repository.UpdateAdminInvoiceRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminInvoiceServiceStruct) UpdateStatusAdminInvoiceService(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error) {
	area, err := s.repository.UpdateStatusAdminInvoiceRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminInvoiceServiceStruct) DeleteAdminInvoiceService(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	area, err := s.repository.DeleteAdminInvoiceRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}
