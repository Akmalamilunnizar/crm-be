package invoice

type InvoiceItem struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required"`
	Qty   int64  `json:"qty" validate:"required,min=1"`
	Total int64  `json:"total" validate:"required"`
}

type CreateAdminInvoiceRequest struct {
	Amount       int64         `json:"amount" validate:"required"`
	CustomerID   string        `json:"customer_id" validate:"required"`
	InvoiceItems []InvoiceItem `json:"invoice_items" validate:"required"`
}

type IdAdminInvoiceRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateStatusAdminInvoiceRequest struct {
	IdAdminInvoiceRequest
	Status string `json:"status" validate:"required,oneof=paid unpaid pending"`
}

type UpdateAdminInvoiceRequest struct {
	IdAdminInvoiceRequest
	CreateAdminInvoiceRequest
}

type PartialPaymentRequest struct {
	IdAdminInvoiceRequest
	Amount int64 `json:"amount" validate:"required,min=1"`
}
