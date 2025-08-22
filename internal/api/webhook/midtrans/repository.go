package midtrans

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type WebhookMidtransRepositoryInterface interface {
	FindWebhookMidtransRepository() ([]entities.Areas, error)
	UpdatePaidWebhookMidtransRepository(id string) (entities.Invoice, error)
	CreateWebhookMidtransRepository(request CreateWebhookMidtransRequest) (entities.Transaction, error)
}

type WebhookMidtransRepositoryStruct struct {
	db *gorm.DB
}

func NewWebhookMidtransRepository(db *gorm.DB) *WebhookMidtransRepositoryStruct {
	return &WebhookMidtransRepositoryStruct{db}
}

func (r WebhookMidtransRepositoryStruct) FindWebhookMidtransRepository() ([]entities.Areas, error) {
	areas := []entities.Areas{}
	tx := r.db.Find(&areas)

	if tx.Error != nil {
		return areas, tx.Error
	}

	return areas, nil
}

func (r WebhookMidtransRepositoryStruct) UpdatePaidWebhookMidtransRepository(id string) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.First(&invoice, "id = ?", id)
	if tx.Error != nil {
		return invoice, tx.Error
	}
	invoice.Status = entities.InvoiceStatusPaid
	tx = r.db.Save(&invoice)
	if tx.Error != nil {
		return invoice, tx.Error
	}

	return invoice, nil
}

func (r WebhookMidtransRepositoryStruct) CreateWebhookMidtransRepository(request CreateWebhookMidtransRequest) (entities.Transaction, error) {
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
