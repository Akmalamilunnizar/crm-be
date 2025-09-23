package recurring_invoice

import (
	"skripsi-be/internal/models/entities"
)

type AdminRecurringInvoiceServiceInterface interface {
	GetAllRecurringInvoices() ([]entities.RecurringInvoice, error)
	GetRecurringInvoiceByID(request IdRecurringInvoiceRequest) (entities.RecurringInvoice, error)
	CreateRecurringInvoice(request CreateRecurringInvoiceRequest, userID string) (entities.RecurringInvoice, error)
	UpdateRecurringInvoice(request UpdateRecurringInvoiceRequest) (entities.RecurringInvoice, error)
	UpdateRecurringInvoiceStatus(request UpdateRecurringInvoiceStatusRequest) (entities.RecurringInvoice, error)
	DeleteRecurringInvoice(request IdRecurringInvoiceRequest) error
	GenerateInvoiceFromRecurring(request GenerateInvoiceRequest) (entities.Invoice, error)
	GetRecurringInvoiceHistory(request IdRecurringInvoiceRequest) ([]entities.RecurringInvoiceHistory, error)
	ProcessDueRecurringInvoices() (int, error)
}

type AdminRecurringInvoiceServiceStruct struct {
	repository AdminRecurringInvoiceRepositoryInterface
}

func NewAdminRecurringInvoiceService(repository AdminRecurringInvoiceRepositoryInterface) *AdminRecurringInvoiceServiceStruct {
	return &AdminRecurringInvoiceServiceStruct{repository}
}

func (s AdminRecurringInvoiceServiceStruct) GetAllRecurringInvoices() ([]entities.RecurringInvoice, error) {
	recurringInvoices, err := s.repository.FindAllRecurringInvoices()
	if err != nil {
		return recurringInvoices, err
	}
	return recurringInvoices, nil
}

func (s AdminRecurringInvoiceServiceStruct) GetRecurringInvoiceByID(request IdRecurringInvoiceRequest) (entities.RecurringInvoice, error) {
	recurringInvoice, err := s.repository.FindRecurringInvoiceByID(request)
	if err != nil {
		return recurringInvoice, err
	}
	return recurringInvoice, nil
}

func (s AdminRecurringInvoiceServiceStruct) CreateRecurringInvoice(request CreateRecurringInvoiceRequest, userID string) (entities.RecurringInvoice, error) {
	recurringInvoice, err := s.repository.CreateRecurringInvoice(request, userID)
	if err != nil {
		return recurringInvoice, err
	}
	return recurringInvoice, nil
}

func (s AdminRecurringInvoiceServiceStruct) UpdateRecurringInvoice(request UpdateRecurringInvoiceRequest) (entities.RecurringInvoice, error) {
	recurringInvoice, err := s.repository.UpdateRecurringInvoice(request)
	if err != nil {
		return recurringInvoice, err
	}
	return recurringInvoice, nil
}

func (s AdminRecurringInvoiceServiceStruct) UpdateRecurringInvoiceStatus(request UpdateRecurringInvoiceStatusRequest) (entities.RecurringInvoice, error) {
	recurringInvoice, err := s.repository.UpdateRecurringInvoiceStatus(request)
	if err != nil {
		return recurringInvoice, err
	}
	return recurringInvoice, nil
}

func (s AdminRecurringInvoiceServiceStruct) DeleteRecurringInvoice(request IdRecurringInvoiceRequest) error {
	err := s.repository.DeleteRecurringInvoice(request)
	if err != nil {
		return err
	}
	return nil
}

func (s AdminRecurringInvoiceServiceStruct) GenerateInvoiceFromRecurring(request GenerateInvoiceRequest) (entities.Invoice, error) {
	invoice, err := s.repository.GenerateInvoiceFromRecurring(request)
	if err != nil {
		return invoice, err
	}
	return invoice, nil
}

func (s AdminRecurringInvoiceServiceStruct) GetRecurringInvoiceHistory(request IdRecurringInvoiceRequest) ([]entities.RecurringInvoiceHistory, error) {
	history, err := s.repository.GetRecurringInvoiceHistory(request)
	if err != nil {
		return history, err
	}
	return history, nil
}

func (s AdminRecurringInvoiceServiceStruct) ProcessDueRecurringInvoices() (int, error) {
	return s.repository.ProcessDueRecurringInvoices()
}
