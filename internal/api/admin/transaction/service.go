package transaction

import (
	"skripsi-be/internal/models/entities"
)

type AdminTransactionServiceInterface interface {
	GetAllAdminTransactionService(request SearchAdminTransactionRequest) ([]entities.Transaction, error)
	GetByIdAdminTransactionService(request IdAdminTransactionRequest) (entities.Transaction, error)
	CreateAdminTransactionService(request CreateAdminTransactionRequest) (entities.Transaction, error)
	UpdateAdminTransactionService(request UpdateAdminTransactionRequest) (entities.Transaction, error)
	DeleteAdminTransactionService(request IdAdminTransactionRequest) (entities.Transaction, error)
	SyncAdminTransactionService() error
}

type AdminTransactionServiceStruct struct {
	repository AdminTransactionRepositoryInterface
}

func NewAdminTransactionService(repository AdminTransactionRepositoryInterface) AdminTransactionServiceStruct {
	return AdminTransactionServiceStruct{repository}
}

func (s AdminTransactionServiceStruct) GetAllAdminTransactionService(request SearchAdminTransactionRequest) ([]entities.Transaction, error) {
	transactions, err := s.repository.FindAdminTransactionRepository(request)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s AdminTransactionServiceStruct) GetByIdAdminTransactionService(request IdAdminTransactionRequest) (entities.Transaction, error) {
	transaction, err := s.repository.FindByIdAdminTransactionRepository(request)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s AdminTransactionServiceStruct) CreateAdminTransactionService(request CreateAdminTransactionRequest) (entities.Transaction, error) {
	transaction, err := s.repository.CreateAdminTransactionRepository(request)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s AdminTransactionServiceStruct) UpdateAdminTransactionService(request UpdateAdminTransactionRequest) (entities.Transaction, error) {
	transaction, err := s.repository.UpdateAdminTransactionRepository(request)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s AdminTransactionServiceStruct) DeleteAdminTransactionService(request IdAdminTransactionRequest) (entities.Transaction, error) {
	transaction, err := s.repository.DeleteAdminTransactionRepository(request)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s AdminTransactionServiceStruct) SyncAdminTransactionService() error {
	err := s.repository.SyncAdminTransactionRepository()
	if err != nil {
		return err
	}

	return nil
}
