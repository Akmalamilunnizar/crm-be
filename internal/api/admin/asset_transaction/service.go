package asset_transaction

import "skripsi-be/internal/models/entities"

type AssetTransactionServiceInterface interface {
	GetAssetTransactionsService(request GetAssetTransactionsRequest) ([]entities.AssetTransaction, error)
	CreateAssetTransactionService(request CreateAssetTransactionRequest) (entities.AssetTransaction, error)
	GetByIdAssetTransactionService(request IdAssetTransactionRequest) (entities.AssetTransaction, error)
	UpdateAssetTransactionService(request UpdateAssetTransactionRequest) (entities.AssetTransaction, error)
	DeleteAssetTransactionService(request IdAssetTransactionRequest) (entities.AssetTransaction, error)
}

type AssetTransactionServiceStruct struct {
	repository AssetTransactionRepositoryInterface
}

func NewAssetTransactionService(repository AssetTransactionRepositoryInterface) AssetTransactionServiceStruct {
	return AssetTransactionServiceStruct{repository}
}

func (s AssetTransactionServiceStruct) GetAssetTransactionsService(request GetAssetTransactionsRequest) ([]entities.AssetTransaction, error) {
	transactions, err := s.repository.FindAssetTransactionsRepository(request)
	if err != nil {
		return transactions, err
	}
	return transactions, err
}

func (s AssetTransactionServiceStruct) GetByIdAssetTransactionService(request IdAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction, err := s.repository.FindByIdAssetTransactionRepository(request)
	if err != nil {
		return transaction, err
	}
	return transaction, err
}

func (s AssetTransactionServiceStruct) CreateAssetTransactionService(request CreateAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction, err := s.repository.CreateAssetTransactionRepository(request)
	if err != nil {
		return transaction, err
	}
	return transaction, err
}

func (s AssetTransactionServiceStruct) UpdateAssetTransactionService(request UpdateAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction, err := s.repository.UpdateAssetTransactionRepository(request)
	if err != nil {
		return transaction, err
	}
	return transaction, err
}

func (s AssetTransactionServiceStruct) DeleteAssetTransactionService(request IdAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction, err := s.repository.DeleteAssetTransactionRepository(request)
	if err != nil {
		return transaction, err
	}
	return transaction, err
}

