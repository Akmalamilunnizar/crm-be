package inventory

import (
	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type InventoryRepositoryInterface interface {
	CreatePurchase(purchase *entities.BarangMasuk, details []entities.DetailBarangMasuk) error
	CreateAssetItems(items []entities.AssetItem) error
	CreateAssetTransaction(transaction *entities.AssetTransaction) error
	CreateTicketAssetTransaction(transaction *entities.TicketAssetTransaction) error
	UpdateAssetItemStatus(assetItemID string, status string) error
	GetAssetItemByID(assetItemID string) (*entities.AssetItem, error)
	GetInventoryStatus(request InventoryStatusRequest) ([]InventoryStatusResponse, error)
	GetAssetByID(assetID string) (*entities.Asset, error)
	GetNextBarangMasukID() (string, error)
}

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepositoryInterface {
	return &InventoryRepository{db: db}
}

// CreatePurchase creates a new purchase record with details
func (r *InventoryRepository) CreatePurchase(purchase *entities.BarangMasuk, details []entities.DetailBarangMasuk) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create main purchase record
	if err := tx.Create(purchase).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create detail records
	for _, detail := range details {
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// CreateAssetItems creates multiple asset items
func (r *InventoryRepository) CreateAssetItems(items []entities.AssetItem) error {
	return r.db.Create(&items).Error
}

// CreateAssetTransaction creates an asset transaction for customer installations
func (r *InventoryRepository) CreateAssetTransaction(transaction *entities.AssetTransaction) error {
	return r.db.Create(transaction).Error
}

// CreateTicketAssetTransaction creates an asset transaction for trouble tickets
func (r *InventoryRepository) CreateTicketAssetTransaction(transaction *entities.TicketAssetTransaction) error {
	return r.db.Create(transaction).Error
}

// UpdateAssetItemStatus updates the status of an asset item
func (r *InventoryRepository) UpdateAssetItemStatus(assetItemID string, status string) error {
	return r.db.Model(&entities.AssetItem{}).
		Where("id = ?", assetItemID).
		Update("status", status).Error
}

// GetAssetItemByID retrieves an asset item by ID
func (r *InventoryRepository) GetAssetItemByID(assetItemID string) (*entities.AssetItem, error) {
	var assetItem entities.AssetItem
	err := r.db.Preload("Asset").First(&assetItem, "id = ?", assetItemID).Error
	if err != nil {
		return nil, err
	}
	return &assetItem, nil
}

// GetInventoryStatus retrieves inventory status with filtering
func (r *InventoryRepository) GetInventoryStatus(request InventoryStatusRequest) ([]InventoryStatusResponse, error) {
	var results []InventoryStatusResponse

	query := r.db.Table("asset_items ai").
		Select(`a.brand, a.model, ai.status, COUNT(*) as count`).
		Joins("JOIN assets a ON ai.asset_id = a.id").
		Group("a.brand, a.model, ai.status")

	// Apply filters
	if request.Brand != nil && *request.Brand != "" {
		query = query.Where("a.brand = ?", *request.Brand)
	}
	if request.Model != nil && *request.Model != "" {
		query = query.Where("a.model = ?", *request.Model)
	}
	if request.Status != nil && *request.Status != "" {
		query = query.Where("ai.status = ?", *request.Status)
	}

	// Execute query to get counts
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result InventoryStatusResponse
		var count int64
		err := rows.Scan(&result.Brand, &result.Model, &result.Status, &count)
		if err != nil {
			return nil, err
		}
		result.Count = count

		// Get individual items for this group
		itemQuery := r.db.Table("asset_items ai").
			Select("ai.id, ai.mac_address, ai.status, ai.created_at, ai.updated_at").
			Joins("JOIN assets a ON ai.asset_id = a.id").
			Where("a.brand = ? AND a.model = ? AND ai.status = ?", result.Brand, result.Model, result.Status)

		var items []AssetItemDetail
		err = itemQuery.Find(&items).Error
		if err != nil {
			return nil, err
		}
		result.Items = items

		results = append(results, result)
	}

	return results, nil
}

// GetAssetByID retrieves an asset by ID
func (r *InventoryRepository) GetAssetByID(assetID string) (*entities.Asset, error) {
	var asset entities.Asset
	err := r.db.First(&asset, "id = ?", assetID).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetNextBarangMasukID generates the next ID for barang masuk
func (r *InventoryRepository) GetNextBarangMasukID() (string, error) {
	var lastID string
	err := r.db.Model(&entities.BarangMasuk{}).
		Select("IdMasuk").
		Order("IdMasuk DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "BM0001", nil // Return first ID if no records exist
	}

	// Extract number from last ID and increment
	if len(lastID) >= 6 {
		// Simple increment logic - in production, you might want more sophisticated ID generation
		return "BM" + lastID[2:], nil
	}

	return "BM0001", nil
}

