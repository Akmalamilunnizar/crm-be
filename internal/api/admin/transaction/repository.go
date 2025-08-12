package transaction

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminTransactionRepositoryInterface interface {
	FindAdminTransactionRepository(request SearchAdminTransactionRequest) ([]entities.Transaction, error)
	FindByIdAdminTransactionRepository(request IdAdminTransactionRequest) (entities.Transaction, error)
	CreateAdminTransactionRepository(request CreateAdminTransactionRequest) (entities.Transaction, error)
	UpdateAdminTransactionRepository(request UpdateAdminTransactionRequest) (entities.Transaction, error)
	DeleteAdminTransactionRepository(request IdAdminTransactionRequest) (entities.Transaction, error)
	SyncAdminTransactionRepository() error
}
type AdminTransactionRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminTransactionRepository(db *gorm.DB) AdminTransactionRepositoryStruct {
	return AdminTransactionRepositoryStruct{db}
}

func (r AdminTransactionRepositoryStruct) FindAdminTransactionRepository(request SearchAdminTransactionRequest) ([]entities.Transaction, error) {
	transactions := []entities.Transaction{}
	tx := r.db
	if request.TypeCash != "" && request.TypeCash != "all" {
		tx = tx.Where("type_cash = ?", request.TypeCash)
	}

	if request.Type != "" && request.Type != "all" {
		tx = tx.Where("type_in_out = ?", request.Type)
	}

	tx = tx.Preload("Account").Find(&transactions)

	if tx.Error != nil {
		return transactions, tx.Error
	}

	return transactions, nil
}
func (r AdminTransactionRepositoryStruct) FindByIdAdminTransactionRepository(request IdAdminTransactionRequest) (entities.Transaction, error) {
	transactions := entities.Transaction{}
	tx := r.db.Preload("Account").First(&transactions, "id = ?", request.Id)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	return transactions, nil
}
func (r AdminTransactionRepositoryStruct) CreateAdminTransactionRepository(request CreateAdminTransactionRequest) (entities.Transaction, error) {
	transactions := entities.Transaction{}
	account := entities.Accounts{}
	copier.Copy(&transactions, &request)

	tx := r.db.Begin()
	txTransaction := r.db.Create(&transactions)
	if txTransaction.Error != nil {
		tx.Rollback()
		return transactions, tx.Error
	}

	txAccount := r.db.First(&account, "id = ?", request.AccountID)
	if txAccount.Error != nil {
		tx.Rollback()
		return transactions, tx.Error
	}

	if condition := transactions.TypeInOut; condition == "credit" {
		account.Saldo = account.Saldo - transactions.Amount
	} else if condition == "debit" {
		account.Saldo = account.Saldo + transactions.Amount
	}

	txAccount = tx.Save(&account)

	if txAccount.Error != nil {
		tx.Rollback()
		return transactions, tx.Error
	}

	tx.Commit()

	return transactions, nil
}

func (r AdminTransactionRepositoryStruct) UpdateAdminTransactionRepository(request UpdateAdminTransactionRequest) (entities.Transaction, error) {
	transactions := entities.Transaction{}
	tx := r.db.First(&transactions, "id = ?", request.Id)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	copier.Copy(&transactions, &request)

	tx = r.db.Save(&transactions)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	return transactions, nil

}
func (r AdminTransactionRepositoryStruct) DeleteAdminTransactionRepository(request IdAdminTransactionRequest) (entities.Transaction, error) {
	transactions := entities.Transaction{}
	tx := r.db.First(&transactions, "id = ?", request.Id)
	if tx.Error != nil {
		return transactions, tx.Error
	}

	tx = r.db.Delete(&transactions)

	return transactions, nil
}

func (r AdminTransactionRepositoryStruct) SyncAdminTransactionRepository() error {
	transactions := []entities.Transaction{}
	tx := r.db.Find(&transactions)
	if tx.Error != nil {
		return tx.Error
	}

	account := entities.Accounts{}
	txAccount := r.db.First(&account, "id = ?", "b82074e7-7acb-40e2-ab33-014f4b09c1f8")
	if txAccount.Error != nil {
		return tx.Error
	}

	var amount int64 = 0
	for _, vtransactions := range transactions {
		if vtransactions.TypeInOut == "credit" {
			amount = amount - vtransactions.Amount
		} else if vtransactions.TypeInOut == "debit" {
			amount = amount + vtransactions.Amount
		}

		account.Saldo = amount
		tx = r.db.Save(&account)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}
