package inventory

// Purchase Stock Request Models
type PurchaseItem struct {
	AssetID      string  `json:"asset_id" validate:"required"`
	SerialNumber *string `json:"serial_number"`
	QtyMasuk     int     `json:"qty_masuk" validate:"required,min=1"`
	HargaSatuan  int     `json:"harga_satuan" validate:"required,min=0"`
	SubTotal     int     `json:"sub_total" validate:"required,min=0"`
}

type CreatePurchaseRequest struct {
	Date  string         `json:"date" validate:"required"`
	Notes *string        `json:"notes"`
	Items []PurchaseItem `json:"items" validate:"required,min=1"`
}

// Asset Deployment Request Models
type CreateDeploymentRequest struct {
	AssetItemID     string  `json:"asset_item_id" validate:"required"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=in out"`
	Notes           *string `json:"notes"`
	// Either customer_installation_id OR trouble_ticket_id must be provided
	CustomerInstallationID *string `json:"customer_installation_id"`
	TroubleTicketID        *uint64 `json:"trouble_ticket_id"`
}

// Inventory Status Query Request Models
type InventoryStatusRequest struct {
	Brand  *string `json:"brand"`
	Model  *string `json:"model"`
	Status *string `json:"status"`
}

// Response Models
type InventoryStatusResponse struct {
	Brand  string            `json:"brand"`
	Model  string            `json:"model"`
	Status string            `json:"status"`
	Count  int64             `json:"count"`
	Items  []AssetItemDetail `json:"items"`
}

type AssetItemDetail struct {
	ID         string `json:"id"`
	MacAddress string `json:"mac_address"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PurchaseResponse struct {
	IdMasuk      string `json:"id_masuk"`
	Date         string `json:"date"`
	TotalItems   int    `json:"total_items"`
	TotalAmount  int    `json:"total_amount"`
	CreatedItems int    `json:"created_items"`
}

type DeploymentResponse struct {
	ID              string `json:"id"`
	AssetItemID     string `json:"asset_item_id"`
	TransactionType string `json:"transaction_type"`
	PreviousStatus  string `json:"previous_status"`
	NewStatus       string `json:"new_status"`
	CreatedAt       string `json:"created_at"`
}

