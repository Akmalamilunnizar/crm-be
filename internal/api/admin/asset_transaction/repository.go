package asset_transaction

import (
	"skripsi-be/internal/models/entities"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AssetTransactionRepositoryInterface interface {
	FindAssetTransactionsRepository(request GetAssetTransactionsRequest) ([]entities.AssetTransaction, error)
	CreateAssetTransactionRepository(request CreateAssetTransactionRequest, createdBy string) (entities.AssetTransaction, error)
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
		if *request.CustomerInstallationID != "" {
			query = query.Where("customer_installation_id = ?", *request.CustomerInstallationID)
		} else {
			// Filter for standalone transactions (empty customer_installation_id)
			query = query.Where("customer_installation_id = ? OR customer_installation_id IS NULL", "")
		}
	} else {
		// If no filter specified, show all transactions including standalone ones
		// This allows fetching both installation-linked and standalone transactions
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

	// Order by transaction date descending (newest first)
	query = query.Order("transaction_date DESC")

	tx := query.Find(&transactions)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	return transactions, nil
}

func (r AssetTransactionRepositoryStruct) CreateAssetTransactionRepository(request CreateAssetTransactionRequest, createdBy string) (entities.AssetTransaction, error) {
	transaction := entities.AssetTransaction{}
	copier.Copy(&transaction, &request)
	
	// Set created_by from context
	transaction.CreatedBy = createdBy
	
	// Handle empty customer_installation_id for standalone goods transactions
	// If empty, we'll still create the transaction but the foreign key constraint might fail
	// For now, we'll use empty string and let the database handle it (if constraint allows NULL or empty)
	if request.CustomerInstallationID == "" {
		// For standalone transactions, we need to set a valid foreign key value
		// Option 1: Use a placeholder value (requires creating a special customer installation record)
		// Option 2: Make the column nullable in database (requires migration)
		// Option 3: Use a default/placeholder installation ID
		// For now, we'll try to create with empty string - if it fails, the error will be caught
		transaction.CustomerInstallationID = "" // Empty string for standalone transactions
	}

	// Parse transaction date if provided
	if request.TransactionDate != "" {
		// Try multiple date formats
		dateFormats := []string{"2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05", "2006-01-02"}
		parsed := false
		for _, format := range dateFormats {
			if parsedDate, err := time.Parse(format, request.TransactionDate); err == nil {
				transaction.TransactionDate = parsedDate
				parsed = true
				break
			}
		}
		if !parsed {
			// If parsing fails, use current time
			transaction.TransactionDate = time.Now()
		}
	} else {
		transaction.TransactionDate = time.Now()
	}

	tx := r.db.Create(&transaction)
	if tx.Error != nil {
		return transaction, tx.Error
	}

	// Reload with relations (skip CustomerInstallation preload if it's empty to avoid errors)
	if transaction.CustomerInstallationID != "" {
		r.db.Preload("Asset").Preload("CustomerInstallation").Preload("User").First(&transaction, "id = ?", transaction.ID)
	} else {
		r.db.Preload("Asset").Preload("User").First(&transaction, "id = ?", transaction.ID)
	}

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
