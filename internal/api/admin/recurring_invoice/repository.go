package recurring_invoice

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
)

type AdminRecurringInvoiceRepositoryInterface interface {
	FindAllRecurringInvoices() ([]entities.RecurringInvoice, error)
	FindRecurringInvoiceByID(request IdRecurringInvoiceRequest) (entities.RecurringInvoice, error)
	CreateRecurringInvoice(request CreateRecurringInvoiceRequest, userID string) (entities.RecurringInvoice, error)
	UpdateRecurringInvoice(request UpdateRecurringInvoiceRequest) (entities.RecurringInvoice, error)
	UpdateRecurringInvoiceStatus(request UpdateRecurringInvoiceStatusRequest) (entities.RecurringInvoice, error)
	DeleteRecurringInvoice(request IdRecurringInvoiceRequest) error
	GenerateInvoiceFromRecurring(request GenerateInvoiceRequest) (entities.Invoice, error)
	GetRecurringInvoiceHistory(request IdRecurringInvoiceRequest) ([]entities.RecurringInvoiceHistory, error)
	ProcessDueRecurringInvoices() (int, error)
}

type AdminRecurringInvoiceRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminRecurringInvoiceRepository(db *gorm.DB) *AdminRecurringInvoiceRepositoryStruct {
	return &AdminRecurringInvoiceRepositoryStruct{db}
}

func (r AdminRecurringInvoiceRepositoryStruct) FindAllRecurringInvoices() ([]entities.RecurringInvoice, error) {
	var recurringInvoices []entities.RecurringInvoice

	tx := r.db.Preload("Customer").
		Preload("CreatedByUser").
		Preload("History.GeneratedInvoice").
		Order("created_at DESC").
		Find(&recurringInvoices)

	if tx.Error != nil {
		return recurringInvoices, tx.Error
	}

	// Parse JSON invoice items for each recurring invoice
	for i := range recurringInvoices {
		if recurringInvoices[i].InvoiceItems != "" {
			var items []entities.RecurringInvoiceItem
			if err := json.Unmarshal([]byte(recurringInvoices[i].InvoiceItems), &items); err == nil {
				recurringInvoices[i].InvoiceItemsData = items
			}
		}
	}

	return recurringInvoices, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) FindRecurringInvoiceByID(request IdRecurringInvoiceRequest) (entities.RecurringInvoice, error) {
	var recurringInvoice entities.RecurringInvoice

	tx := r.db.Preload("Customer").
		Preload("CreatedByUser").
		Preload("History.GeneratedInvoice").
		Where("id = ?", request.Id).
		First(&recurringInvoice)

	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	// Parse JSON invoice items
	if recurringInvoice.InvoiceItems != "" {
		var items []entities.RecurringInvoiceItem
		if err := json.Unmarshal([]byte(recurringInvoice.InvoiceItems), &items); err == nil {
			recurringInvoice.InvoiceItemsData = items
		}
	}

	return recurringInvoice, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) CreateRecurringInvoice(request CreateRecurringInvoiceRequest, userID string) (entities.RecurringInvoice, error) {
	// Calculate next invoice date based on frequency
	nextInvoiceDate := calculateNextInvoiceDate(request.InvoiceDate, request.Frequency)

	// Marshal invoice items to JSON
	itemsJSON, err := json.Marshal(request.InvoiceItems)
	if err != nil {
		return entities.RecurringInvoice{}, fmt.Errorf("failed to marshal invoice items: %v", err)
	}

	recurringInvoice := entities.RecurringInvoice{
		CustomerID:      request.CustomerID,
		Amount:          request.Amount,
		InvoiceDate:     request.InvoiceDate,
		DueDate:         request.DueDate,
		NextInvoiceDate: nextInvoiceDate,
		Frequency:       entities.RecurringInvoiceFrequency(request.Frequency),
		Status:          entities.RecurringInvoiceStatusActive,
		Description:     request.Description,
		InvoiceItems:    string(itemsJSON),
		CreatedBy:       &userID,
	}

	tx := r.db.Create(&recurringInvoice)
	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	// Parse items for response
	var items []entities.RecurringInvoiceItem
	for _, item := range request.InvoiceItems {
		items = append(items, entities.RecurringInvoiceItem{
			Name:  item.Name,
			Price: item.Price,
			Qty:   item.Qty,
			Total: item.Total,
		})
	}
	recurringInvoice.InvoiceItemsData = items

	return recurringInvoice, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) UpdateRecurringInvoice(request UpdateRecurringInvoiceRequest) (entities.RecurringInvoice, error) {
	var recurringInvoice entities.RecurringInvoice

	// Find existing recurring invoice
	tx := r.db.Where("id = ?", request.Id).First(&recurringInvoice)
	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	// Calculate next invoice date based on frequency
	nextInvoiceDate := calculateNextInvoiceDate(request.InvoiceDate, request.Frequency)

	// Marshal invoice items to JSON
	itemsJSON, err := json.Marshal(request.InvoiceItems)
	if err != nil {
		return entities.RecurringInvoice{}, fmt.Errorf("failed to marshal invoice items: %v", err)
	}

	// Update fields
	recurringInvoice.CustomerID = request.CustomerID
	recurringInvoice.Amount = request.Amount
	recurringInvoice.InvoiceDate = request.InvoiceDate
	recurringInvoice.DueDate = request.DueDate
	recurringInvoice.NextInvoiceDate = nextInvoiceDate
	recurringInvoice.Frequency = entities.RecurringInvoiceFrequency(request.Frequency)
	recurringInvoice.Description = request.Description
	recurringInvoice.InvoiceItems = string(itemsJSON)

	tx = r.db.Save(&recurringInvoice)
	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	// Parse items for response
	var items []entities.RecurringInvoiceItem
	for _, item := range request.InvoiceItems {
		items = append(items, entities.RecurringInvoiceItem{
			Name:  item.Name,
			Price: item.Price,
			Qty:   item.Qty,
			Total: item.Total,
		})
	}
	recurringInvoice.InvoiceItemsData = items

	return recurringInvoice, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) UpdateRecurringInvoiceStatus(request UpdateRecurringInvoiceStatusRequest) (entities.RecurringInvoice, error) {
	var recurringInvoice entities.RecurringInvoice

	tx := r.db.Where("id = ?", request.Id).First(&recurringInvoice)
	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	recurringInvoice.Status = entities.RecurringInvoiceStatus(request.Status)

	tx = r.db.Save(&recurringInvoice)
	if tx.Error != nil {
		return recurringInvoice, tx.Error
	}

	// Parse items for response
	if recurringInvoice.InvoiceItems != "" {
		var items []entities.RecurringInvoiceItem
		if err := json.Unmarshal([]byte(recurringInvoice.InvoiceItems), &items); err == nil {
			recurringInvoice.InvoiceItemsData = items
		}
	}

	return recurringInvoice, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) DeleteRecurringInvoice(request IdRecurringInvoiceRequest) error {
	tx := r.db.Where("id = ?", request.Id).Delete(&entities.RecurringInvoice{})
	return tx.Error
}

func (r AdminRecurringInvoiceRepositoryStruct) GenerateInvoiceFromRecurring(request GenerateInvoiceRequest) (entities.Invoice, error) {
	var out entities.Invoice
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// lock the recurring record to avoid double generation
		var recurringInvoice entities.RecurringInvoice
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Customer.Product").First(&recurringInvoice, "id = ?", request.Id).Error; err != nil {
			return err
		}

		// Parse items
		var items []entities.RecurringInvoiceItem
		if err := json.Unmarshal([]byte(recurringInvoice.InvoiceItems), &items); err != nil {
			return fmt.Errorf("failed to parse invoice items: %v", err)
		}

		// Compute dates
		invoiceDate := request.InvoiceDate
		dueDate := request.DueDate
		if invoiceDate == nil {
			invoiceDate = &recurringInvoice.NextInvoiceDate
		}
		if dueDate == nil {
			// preserve original delta
			deltaDays := int(recurringInvoice.DueDate.Sub(recurringInvoice.InvoiceDate).Hours()/24 + 0.5)
			dv := invoiceDate.AddDate(0, 0, deltaDays)
			dueDate = &dv
		}

		// Create invoice
		inv := entities.Invoice{
			CustomerID: recurringInvoice.CustomerID,
			Amount:     recurringInvoice.Amount,
			Link:       fmt.Sprintf("/invoice/%s", recurringInvoice.ID),
			Status:     entities.InvoiceStatusUnpaid,
		}
		if err := tx.Create(&inv).Error; err != nil {
			return err
		}

		// Items
		for _, it := range items {
			name := it.Name
			if name == "" && recurringInvoice.Customer.Product.ID != "" {
				name = recurringInvoice.Customer.Product.Name
			}
			qty := it.Qty
			if qty <= 0 {
				qty = 1
			}
			price := it.Price
			total := it.Total
			if total <= 0 {
				total = price * qty
			}
			ii := entities.InvoiceItems{InvoiceID: inv.ID, Name: name, Price: price, Qty: qty, Total: total}
			// Ensure primary key is set (ID is varchar PK, no BeforeCreate hook)
			ii.ID = uuid.New().String()
			if err := tx.Create(&ii).Error; err != nil {
				return err
			}
		}

		// History
		h := entities.RecurringInvoiceHistory{RecurringInvoiceID: recurringInvoice.ID, GeneratedInvoiceID: inv.ID, InvoiceDate: *invoiceDate, DueDate: *dueDate}
		if err := tx.Create(&h).Error; err != nil {
			return err
		}

		// Update recurring dates (force), with fallback when RowsAffected == 0
		newNext := calculateNextInvoiceDate(*invoiceDate, string(recurringInvoice.Frequency))
		res := tx.Model(&entities.RecurringInvoice{}).Where("id = ?", recurringInvoice.ID).
			Updates(map[string]interface{}{
				"invoice_date":      *invoiceDate,
				"due_date":          *dueDate,
				"next_invoice_date": newNext,
				"updated_at":        time.Now(),
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			if err := tx.Exec(
				"UPDATE recurring_invoices SET invoice_date = ?, due_date = ?, next_invoice_date = ?, updated_at = ? WHERE id = ?",
				*invoiceDate, *dueDate, newNext, time.Now(), recurringInvoice.ID,
			).Error; err != nil {
				return err
			}
		}

		out = inv
		return nil
	})
	if err != nil {
		return entities.Invoice{}, err
	}
	// reload invoice with relations for response (explicit condition to avoid malformed WHERE)
	r.db.Preload("Customer").Preload("InvoiceItems").First(&out, "id = ?", out.ID)
	return out, nil
}

func (r AdminRecurringInvoiceRepositoryStruct) GetRecurringInvoiceHistory(request IdRecurringInvoiceRequest) ([]entities.RecurringInvoiceHistory, error) {
	var history []entities.RecurringInvoiceHistory

	tx := r.db.Preload("GeneratedInvoice").
		Where("recurring_invoice_id = ?", request.Id).
		Order("generated_at DESC").
		Find(&history)

	return history, tx.Error
}

// Helper function to calculate next invoice date
func calculateNextInvoiceDate(currentDate time.Time, frequency string) time.Time {
	switch frequency {
	case "monthly":
		return currentDate.AddDate(0, 1, 0)
	case "quarterly":
		return currentDate.AddDate(0, 3, 0)
	case "yearly":
		return currentDate.AddDate(1, 0, 0)
	default:
		return currentDate.AddDate(0, 1, 0) // Default to monthly
	}
}

// ProcessDueRecurringInvoices finds all active recurring invoices with next_invoice_date <= now
// and generates invoices for them. Returns the number of invoices generated.
func (r AdminRecurringInvoiceRepositoryStruct) ProcessDueRecurringInvoices() (int, error) {
	var due []entities.RecurringInvoice
	tx := r.db.Where("status = ? AND next_invoice_date <= ?", entities.RecurringInvoiceStatusActive, time.Now()).Find(&due)
	if tx.Error != nil {
		return 0, tx.Error
	}

	generated := 0
	for _, rec := range due {
		_, err := r.GenerateInvoiceFromRecurring(GenerateInvoiceRequest{IdRecurringInvoiceRequest: IdRecurringInvoiceRequest{Id: rec.ID}})
		if err == nil {
			generated++
		}
	}
	return generated, nil
}
