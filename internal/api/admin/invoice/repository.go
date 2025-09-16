package invoice

import (
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"skripsi-be/internal/models/entities"
)

type AdminInvoiceRepositoryInterface interface {
	FindAdminInvoiceRepository() ([]entities.Invoice, error)
	CreateAdminInvoiceRepository(request CreateAdminInvoiceRequest) (entities.Invoice, error)
	FindByIdAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
	UpdateAdminInvoiceRepository(request UpdateAdminInvoiceRequest) (entities.Invoice, error)
	UpdateStatusAdminInvoiceRepository(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error)
	DeleteAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
	ProcessPartialPaymentRepository(request PartialPaymentRequest) (entities.Invoice, error)
	MarkPdfViewedRepository(request IdAdminInvoiceRequest) (entities.Invoice, error)
}

type AdminInvoiceRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminInvoiceRepository(db *gorm.DB) *AdminInvoiceRepositoryStruct {
	return &AdminInvoiceRepositoryStruct{db}
}

func (r AdminInvoiceRepositoryStruct) FindAdminInvoiceRepository() ([]entities.Invoice, error) {
	invoices := []entities.Invoice{}
	tx := r.db.Preload("Customer.Product").Preload("InvoiceItems.Invoice").Order("createdAt DESC").Find(&invoices)

	if tx.Error != nil {
		return invoices, tx.Error
	}

	// For each invoice, calculate total paid from transactions
	for i := range invoices {
		var totalPaid int64
		r.db.Model(&entities.Transaction{}).Where("invoice_id = ?", invoices[i].ID).Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)

		// Create a virtual transaction object to hold total paid
		if totalPaid > 0 {
			invoices[i].Transaction = entities.Transaction{
				Amount: totalPaid,
			}
		}
	}

	return invoices, nil
}

func (r AdminInvoiceRepositoryStruct) FindByIdAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.Preload("Customer.Product").Preload("InvoiceItems.Invoice").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Calculate total amount from invoice items if not set or incorrect
	if len(invoice.InvoiceItems) > 0 {
		var totalAmount int64 = 0
		for _, item := range invoice.InvoiceItems {
			totalAmount += item.Total
		}

		// Update invoice amount if it's different from calculated total
		if invoice.Amount != totalAmount {
			invoice.Amount = totalAmount
			// Update in database
			r.db.Model(&invoice).Update("amount", totalAmount)
		}
	}

	// Calculate total paid from all transactions for this invoice
	var totalPaid int64
	r.db.Model(&entities.Transaction{}).Where("invoice_id = ?", request.Id).Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)

	// Create a virtual transaction object to hold total paid
	if totalPaid > 0 {
		invoice.Transaction = entities.Transaction{
			Amount: totalPaid,
		}
	}

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) CreateAdminInvoiceRepository(request CreateAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	copier.Copy(&invoice, &request)
	tx := r.db.Begin()
	txInvoice := tx.Model(&invoice).Create(&invoice)
	if txInvoice.Error != nil {
		tx.Rollback()
		return entities.Invoice{}, tx.Error
	}

	// Calculate total amount from all invoice items and save them
	var totalAmount int64 = 0
	for i := range invoice.InvoiceItems {
		invoice.InvoiceItems[i].InvoiceID = invoice.ID
		invoice.InvoiceItems[i].Total = invoice.InvoiceItems[i].Price * invoice.InvoiceItems[i].Qty
		totalAmount += invoice.InvoiceItems[i].Total

		// Save each invoice item
		txInvoiceItem := tx.Create(&invoice.InvoiceItems[i])
		if txInvoiceItem.Error != nil {
			tx.Rollback()
			return entities.Invoice{}, txInvoiceItem.Error
		}
	}

	// Set total amount to sum of all items
	invoice.Amount = totalAmount
	txInvoice = tx.Model(&invoice).Updates(entities.Invoice{Amount: invoice.Amount})
	if txInvoice.Error != nil {
		tx.Rollback()
		return entities.Invoice{}, tx.Error
	}

	tx.Commit()

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) UpdateAdminInvoiceRepository(request UpdateAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}
	copier.Copy(&invoice, &request)

	r.db.Save(&invoice)
	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) UpdateStatusAdminInvoiceRepository(request UpdateStatusAdminInvoiceRequest) (entities.Invoice, error) {
	invoice := entities.Invoice{}
	tx := r.db.First(&invoice, "id = ?", request.Id)

	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Update only the status field directly
	tx = r.db.Model(&invoice).Update("status", request.Status)
	if tx.Error != nil {
		return invoice, tx.Error
	}

	// Reload the invoice to get updated data
	r.db.Preload("Customer.Product").Preload("InvoiceItems.Invoice").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) DeleteAdminInvoiceRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {

	invoice := entities.Invoice{}

	tx := r.db.First(&invoice, "id = ?", request.Id)
	if tx.Error != nil {
		return invoice, tx.Error
	}

	tx = r.db.Delete(&invoice)
	if tx.Error != nil {
		return invoice, tx.Error
	}
	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) ProcessPartialPaymentRepository(request PartialPaymentRequest) (entities.Invoice, error) {
	// Start transaction
	tx := r.db.Begin()

	// Get invoice with related data
	invoice := entities.Invoice{}
	err := tx.Preload("Customer").Preload("InvoiceItems").First(&invoice, "id = ?", request.Id).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Validate invoice status (should be pending or unpaid)
	if invoice.Status == entities.InvoiceStatusPaid {
		tx.Rollback()
		return invoice, fmt.Errorf("invoice is already fully paid")
	}

	// Calculate current total paid from existing transactions
	var currentTotalPaid int64
	tx.Model(&entities.Transaction{}).Where("invoice_id = ?", request.Id).Select("COALESCE(SUM(amount), 0)").Scan(&currentTotalPaid)

	// Calculate new total paid
	newTotalPaid := currentTotalPaid + request.Amount

	// Validate payment amount doesn't exceed invoice amount
	if newTotalPaid > invoice.Amount {
		tx.Rollback()
		return invoice, fmt.Errorf("payment amount exceeds outstanding balance")
	}

	// Create payment transaction
	transaction := entities.Transaction{
		AccountID:   "b82074e7-7acb-40e2-ab33-014f4b09c1f8", // Default account ID
		InvoiceID:   request.Id,
		TypeCash:    entities.TransactionsTypeCashInternet,
		TypeInOut:   entities.TransactionsTypeInOutIn,
		Description: fmt.Sprintf("Partial payment for invoice %s", request.Id),
		Amount:      request.Amount,
		Category:    "partial_payment",
		Method:      "manual",
	}

	err = tx.Create(&transaction).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Update invoice status based on payment
	var newStatus entities.InvoiceStatus
	if newTotalPaid >= invoice.Amount {
		newStatus = entities.InvoiceStatusPaid
	} else {
		newStatus = entities.InvoiceStatusPending
	}

	err = tx.Model(&invoice).Update("status", newStatus).Error
	if err != nil {
		tx.Rollback()
		return invoice, err
	}

	// Commit transaction
	tx.Commit()

	// Reload invoice with updated data
	r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	return invoice, nil
}

func (r AdminInvoiceRepositoryStruct) MarkPdfViewedRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {
	var invoice entities.Invoice

	// Check if invoice exists
	err := r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id).Error
	if err != nil {
		return invoice, err
	}

	// Check if PDF has already been viewed
	if invoice.PdfViewed {
		return invoice, fmt.Errorf("PDF has already been viewed and cannot be accessed again")
	}

	// Mark PDF as viewed with current timestamp
	now := time.Now()
	err = r.db.Model(&invoice).Updates(map[string]interface{}{
		"pdf_viewed":    true,
		"pdf_viewed_at": &now,
	}).Error

	if err != nil {
		return invoice, err
	}

	// Reload invoice with updated data
	r.db.Preload("Customer").Preload("InvoiceItems").Preload("Transaction").First(&invoice, "id = ?", request.Id)

	return invoice, nil
}
