package asset_transaction

type CreateAssetTransactionRequest struct {
	CustomerInstallationID string  `json:"customer_installation_id"` // Optional - empty string for standalone goods transactions
	AssetID                string  `json:"asset_id" validate:"required"`
	TransactionType        string  `json:"transaction_type" validate:"required,oneof=out in"`
	Quantity               int     `json:"quantity" validate:"required,min=1"`
	Notes                  *string `json:"notes"`
	TransactionDate        string  `json:"transaction_date"`
}

type IdAssetTransactionRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAssetTransactionRequest struct {
	IdAssetTransactionRequest
	CreateAssetTransactionRequest
}

type GetAssetTransactionsRequest struct {
	AssetID                *string `json:"asset_id"`
	CustomerInstallationID *string `json:"customer_installation_id"`
	TransactionType        *string `json:"transaction_type"`
	DateFrom               *string `json:"date_from"`
	DateTo                 *string `json:"date_to"`
}
