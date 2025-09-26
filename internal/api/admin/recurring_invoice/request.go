package recurring_invoice

import "time"

type RecurringInvoiceItem struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,min=1"`
	Qty   int64  `json:"qty" validate:"required,min=1"`
	Total int64  `json:"total" validate:"required,min=1"`
}

type CreateRecurringInvoiceRequest struct {
	CustomerID   string                 `json:"customer_id" validate:"required"`
	Amount       int64                  `json:"amount" validate:"required,min=1"`
	InvoiceDate  time.Time              `json:"invoice_date" validate:"required"`
	DueDate      time.Time              `json:"due_date" validate:"required"`
	Frequency    string                 `json:"frequency" validate:"required,oneof=monthly quarterly yearly"`
	Description  *string                `json:"description"`
	InvoiceItems []RecurringInvoiceItem `json:"invoice_items" validate:"required,min=1"`
}

type IdRecurringInvoiceRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateRecurringInvoiceRequest struct {
	IdRecurringInvoiceRequest
	CreateRecurringInvoiceRequest
}

type UpdateRecurringInvoiceStatusRequest struct {
	IdRecurringInvoiceRequest
	Status string `json:"status" validate:"required,oneof=active stopped completed"`
}

type GenerateInvoiceRequest struct {
	IdRecurringInvoiceRequest
	InvoiceDate *time.Time `json:"invoice_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
