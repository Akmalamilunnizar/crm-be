package asset_transaction

import (
	"skripsi-be/internal/models/entities"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AssetTransactionRepositoryInterface interface {
	FindAssetTransactionsRepository(request GetAssetTransactionsRequest) ([]entities.AssetTransaction, error)
	CreateAssetTransactionRepository(request CreateAssetTransactionRequest) (entities.AssetTransaction, error)
	FindByIdAssetTransactionRepository(request IdAssetTransactionRequest) (entities.AssetTransaction, error)
	UpdateAssetTransactionRepository(request UpdateAssetTransactionRequest) (entities.AssetTransaction, error)
	DeleteAssetTransactionRepository(request IdAssetTransactionRequest) (entities.AssetTransaction, error)
}

type AssetTransactionRepositoryStruct struct {
	db *gorm.DB
}

func NewAssetTransactionRepository(db *gorm.DB) *AssetTransactionRepositoryStruct {
	return &AssetTransactionRepositoryStruct{db}
}

func (r AssetTransactionRepositoryStruct) FindAssetTransactionsRepository(request GetAssetTransactionsRequest) ([]entities.AssetTransaction, error) {
	transactions := []entities.AssetTransaction{}
	query := r.db.Preload("Asset").Preload("CustomerInstallation").Preload("User")

	if request.AssetID != nil {
		query = query.Where("asset_id = ?", *request.AssetID)
	}
	if request.CustomerInstallationID != nil {
		query = query.Where("customer_installation_id = ?", *request.CustomerInstallationID)
	}
	if request.TransactionType != nil {
		query = query.Where("transaction_type = ?", *request.TransactionType)
	}
	if request.DateFrom != nil {
		query = query.Where("transaction_date >= ?", *request.DateFrom)
	}
	if request.DateTo != nil {
		query = query.Where("transaction_date <= ?", *request.DateTo)
	}

	tx := query.Find(&transactions)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	return transactions, nil
}

func (r AssetTransactionRepositoryStruct) CreateAssetTransactionRepository(request CreateAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction := entities.AssetTransaction{}
	copier.Copy(&transaction, &request)

	// Parse transaction date if provided
	if request.TransactionDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", request.TransactionDate); err == nil {
			transaction.TransactionDate = parsedDate
		}
	}

	tx := r.db.Create(&transaction)
	if tx.Error != nil {
		return transaction, tx.Error
	}

	// Reload with relations
	r.db.Preload("Asset").Preload("CustomerInstallation").Preload("User").First(&transaction, "id = ?", transaction.ID)

	return transaction, nil
}

func (r AssetTransactionRepositoryStruct) FindByIdAssetTransactionRepository(request IdAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction := entities.AssetTransaction{}
	tx := r.db.Preload("Asset").Preload("CustomerInstallation").Preload("User").First(&transaction, "id = ?", request.Id)
	if tx.Error != nil {
		return transaction, tx.Error
	}
	return transaction, nil
}

func (r AssetTransactionRepositoryStruct) UpdateAssetTransactionRepository(request UpdateAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction := entities.AssetTransaction{}
	tx := r.db.First(&transaction, "id = ?", request.Id)
	if tx.Error != nil {
		return transaction, tx.Error
	}
	copier.Copy(&transaction, &request.CreateAssetTransactionRequest)

	// Parse transaction date if provided
	if request.CreateAssetTransactionRequest.TransactionDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", request.CreateAssetTransactionRequest.TransactionDate); err == nil {
			transaction.TransactionDate = parsedDate
		}
	}

	tx = r.db.Save(&transaction)
	if tx.Error != nil {
		return transaction, tx.Error
	}

	// Reload with relations
	r.db.Preload("Asset").Preload("CustomerInstallation").Preload("User").First(&transaction, "id = ?", transaction.ID)

	return transaction, nil
}

func (r AssetTransactionRepositoryStruct) DeleteAssetTransactionRepository(request IdAssetTransactionRequest) (entities.AssetTransaction, error) {
	transaction := entities.AssetTransaction{}
	tx := r.db.First(&transaction, "id = ?", request.Id)
	if tx.Error != nil {
		return transaction, tx.Error
	}

	tx = r.db.Delete(&transaction)
	if tx.Error != nil {
		return transaction, tx.Error
	}

	return transaction, nil
}

