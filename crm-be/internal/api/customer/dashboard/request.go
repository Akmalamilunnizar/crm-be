package dashboard

type CreatePaymentCustomerDashboardRequest struct {
	InvoiceId string `json:"invoice_id" validate:"required"`
}

type idUser struct {
	UserId string `json:"user_id" `
}

type SearchInvoice struct {
	idUser
	InvoiceId string `json:"invoice_id" `
}
