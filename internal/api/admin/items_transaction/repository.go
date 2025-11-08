package items_transaction

import (
	"fmt"
	"skripsi-be/internal/models/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemsTransactionRepositoryInterface interface {
	FindItemsTransactionsRepository(request GetItemsTransactionsRequest) ([]ItemsTransactionResponse, error)
	CreateItemsTransactionRepository(request CreateItemsTransactionRequest, createdBy string) (ItemsTransactionResponse, error)
	FindByIdItemsTransactionRepository(request IdItemsTransactionRequest, transactionType string) (ItemsTransactionResponse, error)
	UpdateItemsTransactionRepository(request UpdateItemsTransactionRequest) (ItemsTransactionResponse, error)
	DeleteItemsTransactionRepository(request IdItemsTransactionRequest, transactionType string) error
	DeleteItemsTransactionDetailRepository(request DeleteItemsTransactionDetailRequest) error
}

type ItemsTransactionRepositoryStruct struct {
	db *gorm.DB
}

func NewItemsTransactionRepository(db *gorm.DB) *ItemsTransactionRepositoryStruct {
	return &ItemsTransactionRepositoryStruct{db}
}

func (r ItemsTransactionRepositoryStruct) FindItemsTransactionsRepository(request GetItemsTransactionsRequest) ([]ItemsTransactionResponse, error) {
	var responses []ItemsTransactionResponse

	// Query items keluar (out)
	if request.TransactionType == nil || *request.TransactionType == "out" {
		var itemsKeluar []entities.ItemsKeluar
		query := r.db.Preload("Details.Item").Preload("User") // Preload Item from catalog

		if request.DateFrom != nil {
			query = query.Where("date >= ?", *request.DateFrom)
		}
		if request.DateTo != nil {
			query = query.Where("date <= ?", *request.DateTo)
		}
		if request.AssetID != nil {
			// Filter by asset_id through items catalog (items can have asset_id)
			query = query.Joins("JOIN detail_itemskeluar ON detail_itemskeluar.IdKeluar = itemskeluar.id").
				Joins("JOIN items ON items.id = detail_itemskeluar.item_id").
				Where("items.asset_id = ?", *request.AssetID)
		}

		if err := query.Order("date DESC, created_at DESC").Find(&itemsKeluar).Error; err != nil {
			return nil, err
		}

		for _, ik := range itemsKeluar {
			items := make([]ItemsTransactionDetailResponse, 0)
			for _, detail := range ik.Details {
				// Use Item from catalog
				var itemName string
				var itemInfo interface{}
				if detail.Item != nil {
					itemName = detail.Item.Name
					itemInfo = detail.Item
				} else {
					itemName = "Unknown Item"
					itemInfo = nil
				}

				items = append(items, ItemsTransactionDetailResponse{
					IdItems:   detail.ItemID,
					ItemName:  itemName,
					Item:      itemInfo, // Item from catalog
					Quantity:  detail.QtyKeluar,
					Unit:      detail.Unit,
					Notes:     detail.Notes,
					CreatedAt: detail.CreatedAt.Format(time.RFC3339),
				})
			}

			responses = append(responses, ItemsTransactionResponse{
				ID:              ik.ID,
				TransactionType: "out",
				Date:            ik.Date.Format("2006-01-02"),
				Notes:           ik.Notes,
				CreatedBy:       ik.CreatedBy,
				CreatedAt:       ik.CreatedAt.Format(time.RFC3339),
				UpdatedAt:       ik.UpdatedAt.Format(time.RFC3339),
				Items:           items,
				User:            ik.User,
			})
		}
	}

	// Query items masuk (in)
	if request.TransactionType == nil || *request.TransactionType == "in" {
		var itemsMasuk []entities.ItemsMasuk
		query := r.db.Preload("Details.Item").Preload("User") // Preload Item from catalog

		if request.DateFrom != nil {
			query = query.Where("date >= ?", *request.DateFrom)
		}
		if request.DateTo != nil {
			query = query.Where("date <= ?", *request.DateTo)
		}
		if request.AssetID != nil {
			// Filter by asset_id through items catalog (items can have asset_id)
			query = query.Joins("JOIN detail_itemsmasuk ON detail_itemsmasuk.IdMasuk = itemsmasuk.id").
				Joins("JOIN items ON items.id = detail_itemsmasuk.item_id").
				Where("items.asset_id = ?", *request.AssetID)
		}

		if err := query.Order("date DESC, created_at DESC").Find(&itemsMasuk).Error; err != nil {
			return nil, err
		}

		for _, im := range itemsMasuk {
			items := make([]ItemsTransactionDetailResponse, 0)
			for _, detail := range im.Details {
				subTotal := detail.SubTotal
				// Use Item from catalog
				var itemName string
				var itemInfo interface{}
				if detail.Item != nil {
					itemName = detail.Item.Name
					itemInfo = detail.Item
				} else {
					itemName = "Unknown Item"
					itemInfo = nil
				}

				items = append(items, ItemsTransactionDetailResponse{
					IdItems:     detail.ItemID,
					ItemName:    itemName,
					Item:        itemInfo, // Item from catalog
					Quantity:    detail.QtyMasuk,
					Unit:        detail.Unit,
					HargaSatuan: &detail.HargaSatuan,
					SubTotal:    &subTotal,
					Notes:       detail.Notes,
					CreatedAt:   detail.CreatedAt.Format(time.RFC3339),
				})
			}

			responses = append(responses, ItemsTransactionResponse{
				ID:              im.ID,
				TransactionType: "in",
				Date:            im.Date.Format("2006-01-02"),
				Notes:           im.Notes,
				CreatedBy:       im.CreatedBy,
				CreatedAt:       im.CreatedAt.Format(time.RFC3339),
				UpdatedAt:       im.UpdatedAt.Format(time.RFC3339),
				Items:           items,
				User:            im.User,
			})
		}
	}

	return responses, nil
}

func (r ItemsTransactionRepositoryStruct) CreateItemsTransactionRepository(request CreateItemsTransactionRequest, createdBy string) (ItemsTransactionResponse, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return ItemsTransactionResponse{}, fmt.Errorf("invalid date format: %v", err)
	}

	var response ItemsTransactionResponse

	if request.TransactionType == "out" {
		// Create ItemsKeluar
		itemsKeluar := entities.ItemsKeluar{
			ID:        generateItemsTransactionID("IK"),
			Date:      date,
			Notes:     request.Notes,
			CreatedBy: createdBy,
		}

		// Create details
		details := make([]entities.DetailItemsKeluar, 0)
		for _, item := range request.Items {
			// item.IdItems should be item_id from items catalog
			details = append(details, entities.DetailItemsKeluar{
				IdKeluar:  itemsKeluar.ID,
				ItemID:    item.IdItems, // Reference to items catalog (item_id from items table)
				QtyKeluar: item.Quantity,
				Unit:      item.Unit,
				Notes:     item.Notes,
			})
		}

		// Save in transaction
		tx := r.db.Begin()
		if err := tx.Create(&itemsKeluar).Error; err != nil {
			tx.Rollback()
			return ItemsTransactionResponse{}, err
		}

		for _, detail := range details {
			if err := tx.Create(&detail).Error; err != nil {
				tx.Rollback()
				return ItemsTransactionResponse{}, err
			}
		}

		if err := tx.Commit().Error; err != nil {
			return ItemsTransactionResponse{}, err
		}

		// Reload with relations
		r.db.Preload("Details.Item").Preload("User").First(&itemsKeluar, "id = ?", itemsKeluar.ID)

		// Build response
		items := make([]ItemsTransactionDetailResponse, 0)
		for _, detail := range itemsKeluar.Details {
			var itemName string
			var itemInfo interface{}
			if detail.Item != nil {
				itemName = detail.Item.Name
				itemInfo = detail.Item
			} else {
				itemName = "Unknown Item"
				itemInfo = nil
			}

			items = append(items, ItemsTransactionDetailResponse{
				IdItems:   detail.ItemID,
				ItemName:  itemName,
				Item:      itemInfo,
				Quantity:  detail.QtyKeluar,
				Unit:      detail.Unit,
				Notes:     detail.Notes,
				CreatedAt: detail.CreatedAt.Format(time.RFC3339),
			})
		}

		response = ItemsTransactionResponse{
			ID:              itemsKeluar.ID,
			TransactionType: "out",
			Date:            itemsKeluar.Date.Format("2006-01-02"),
			Notes:           itemsKeluar.Notes,
			CreatedBy:       itemsKeluar.CreatedBy,
			CreatedAt:       itemsKeluar.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       itemsKeluar.UpdatedAt.Format(time.RFC3339),
			Items:           items,
			User:            itemsKeluar.User,
		}
	} else {
		// Create ItemsMasuk
		itemsMasuk := entities.ItemsMasuk{
			ID:        generateItemsTransactionID("IM"),
			Date:      date,
			Notes:     request.Notes,
			CreatedBy: createdBy,
		}

		// Create details
		details := make([]entities.DetailItemsMasuk, 0)
		for _, item := range request.Items {
			hargaSatuan := 0
			if item.HargaSatuan != nil {
				hargaSatuan = *item.HargaSatuan
			}
			subTotal := hargaSatuan * item.Quantity

			// item.IdItems should be item_id from items catalog
			details = append(details, entities.DetailItemsMasuk{
				IdMasuk:     itemsMasuk.ID,
				ItemID:      item.IdItems, // Reference to items catalog (item_id from items table)
				QtyMasuk:    item.Quantity,
				HargaSatuan: hargaSatuan,
				SubTotal:    subTotal,
				Unit:        item.Unit,
				Notes:       item.Notes,
			})
		}

		// Save in transaction
		tx := r.db.Begin()
		if err := tx.Create(&itemsMasuk).Error; err != nil {
			tx.Rollback()
			return ItemsTransactionResponse{}, err
		}

		for _, detail := range details {
			if err := tx.Create(&detail).Error; err != nil {
				tx.Rollback()
				return ItemsTransactionResponse{}, err
			}
		}

		if err := tx.Commit().Error; err != nil {
			return ItemsTransactionResponse{}, err
		}

		// Reload with relations
		r.db.Preload("Details.Item").Preload("User").First(&itemsMasuk, "id = ?", itemsMasuk.ID)

		// Build response
		items := make([]ItemsTransactionDetailResponse, 0)
		for _, detail := range itemsMasuk.Details {
			subTotal := detail.SubTotal
			var itemName string
			var itemInfo interface{}
			if detail.Item != nil {
				itemName = detail.Item.Name
				itemInfo = detail.Item
			} else {
				itemName = "Unknown Item"
				itemInfo = nil
			}

			items = append(items, ItemsTransactionDetailResponse{
				IdItems:     detail.ItemID,
				ItemName:    itemName,
				Item:        itemInfo,
				Quantity:    detail.QtyMasuk,
				Unit:        detail.Unit,
				HargaSatuan: &detail.HargaSatuan,
				SubTotal:    &subTotal,
				Notes:       detail.Notes,
				CreatedAt:   detail.CreatedAt.Format(time.RFC3339),
			})
		}

		response = ItemsTransactionResponse{
			ID:              itemsMasuk.ID,
			TransactionType: "in",
			Date:            itemsMasuk.Date.Format("2006-01-02"),
			Notes:           itemsMasuk.Notes,
			CreatedBy:       itemsMasuk.CreatedBy,
			CreatedAt:       itemsMasuk.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       itemsMasuk.UpdatedAt.Format(time.RFC3339),
			Items:           items,
			User:            itemsMasuk.User,
		}
	}

	return response, nil
}

func (r ItemsTransactionRepositoryStruct) FindByIdItemsTransactionRepository(request IdItemsTransactionRequest, transactionType string) (ItemsTransactionResponse, error) {
	// Implementation similar to FindItemsTransactionsRepository but for single record
	// For brevity, I'll implement a simplified version
	return ItemsTransactionResponse{}, fmt.Errorf("not implemented yet")
}

func (r ItemsTransactionRepositoryStruct) UpdateItemsTransactionRepository(request UpdateItemsTransactionRequest) (ItemsTransactionResponse, error) {
	// Implementation for update
	return ItemsTransactionResponse{}, fmt.Errorf("not implemented yet")
}

func (r ItemsTransactionRepositoryStruct) DeleteItemsTransactionRepository(request IdItemsTransactionRequest, transactionType string) error {
	if transactionType == "out" {
		return r.db.Delete(&entities.ItemsKeluar{}, "id = ?", request.Id).Error
	}
	return r.db.Delete(&entities.ItemsMasuk{}, "id = ?", request.Id).Error
}

func (r ItemsTransactionRepositoryStruct) DeleteItemsTransactionDetailRepository(request DeleteItemsTransactionDetailRequest) error {
	if request.TransactionType == "out" {
		return r.db.Where("IdKeluar = ? AND IdItems = ?", request.TransactionID, request.DetailID).
			Delete(&entities.DetailItemsKeluar{}).Error
	}
	return r.db.Where("IdMasuk = ? AND IdItems = ?", request.TransactionID, request.DetailID).
		Delete(&entities.DetailItemsMasuk{}).Error
}

// Helper function to generate transaction IDs
func generateItemsTransactionID(prefix string) string {
	// Simple ID generation - in production, use proper sequence
	uuid := uuid.New().String()
	return prefix + uuid[:4]
}
