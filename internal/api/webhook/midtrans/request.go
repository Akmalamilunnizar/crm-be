package midtrans

type CreateWebhookMootaRequest struct {
	Account_number int64  `json:"account_number"`
	Date           string `json:"date"`
	Description    string `json:"description"`
	Amount         int64  `json:"amount"`
	Type           string `json:"type"`
	Balance        int64  `json:"balance"`
	Updated_at     string `json:"updated_at"`
	Created_at     string `json:"created_at"`
	Mutation_id    string `json:"mutation_id"`
	Token          string `json:"token"`
	Bank_id        string `json:"bank_id"`
}
type MidtransNotification struct {
	TransactionType   string `json:"transaction_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	Issuer            string `json:"issuer"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	ExpiryTime        string `json:"expiry_time"`
	Currency          string `json:"currency"`
	Acquirer          string `json:"acquirer"`
}

type CreateWebhookMidtransRequest struct {
	AccountID   string `json:"account_id" validate:"required"`
	TypeCash    string `json:"type_cash" validate:"required"`
	TypeInOut   string `json:"type_in_out" validate:"required"`
	Description string `json:"description"`
	Amount      int64  `json:"amount" validate:"required"`
}
type IdWebhookMootaRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateWebhookMootaRequest struct {
	IdWebhookMootaRequest
	CreateWebhookMootaRequest
}
