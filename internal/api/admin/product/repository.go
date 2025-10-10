package product

import (
	"fmt"
	"log"
	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type AdminProductRepositoryInterface interface {
	FindAdminProductRepository() []entities.Products
}

type AdminProductRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminProductRepository(db *gorm.DB) AdminProductRepositoryStruct {

	return AdminProductRepositoryStruct{db}
}

func (r AdminProductRepositoryStruct) FindAdminProductRepository() ([]entities.Products, error) {
	products := []entities.Products{}

	// Debug: Log the database connection info
	log.Printf("🔍 DEBUG: Fetching products from database...")

	tx := r.db.Find(&products)
	if tx.Error != nil {
		log.Printf("❌ ERROR: Failed to fetch products: %v", tx.Error)
		return nil, tx.Error
	}

	// Debug: Log the fetched products
	log.Printf("✅ DEBUG: Fetched %d products", len(products))
	for i, product := range products {
		var downloadSpeed, uploadSpeed string
		if product.DownloadSpeedMbps != nil {
			downloadSpeed = fmt.Sprintf("%d", *product.DownloadSpeedMbps)
		} else {
			downloadSpeed = "nil"
		}
		if product.UploadSpeedMbps != nil {
			uploadSpeed = fmt.Sprintf("%d", *product.UploadSpeedMbps)
		} else {
			uploadSpeed = "nil"
		}
		log.Printf("📦 Product %d: ID=%s, Name=%s, DownloadSpeed=%s, UploadSpeed=%s",
			i+1, product.ID, product.Name, downloadSpeed, uploadSpeed)
	}

	return products, nil
}

func (r AdminProductRepositoryStruct) FindByIdAdminProductRepository(id IdAdminProductRequest) (entities.Products, error) {
	product := entities.Products{}
	tx := r.db.First(&product, "id = ?", id.Id)
	if tx.Error != nil {
		return entities.Products{}, tx.Error
	}
	return product, nil
}

func (r AdminProductRepositoryStruct) CreateAdminProductRepository(request CreateAdminProductRequest) (entities.Products, error) {
	product := entities.Products{}

	product.Name = request.Name
	product.Price = request.Price
	product.Description = request.Description

	tx := r.db.Create(&product)

	if tx.Error != nil {
		return entities.Products{}, tx.Error
	}

	return product, nil
}

func (r AdminProductRepositoryStruct) UpdateAdminProductRepository(request UpdateAdminProductRequest) (entities.Products, error) {
	product := entities.Products{}

	tx := r.db.First(&product, "id = ? ", request.Id)
	if tx.Error != nil {
		return entities.Products{}, tx.Error
	}

	product.Name = request.Name
	product.Price = request.Price
	product.Description = request.Description

	tx = r.db.Save(&product)

	if tx.Error != nil {
		return entities.Products{}, tx.Error
	}

	return product, nil
}

func (r AdminProductRepositoryStruct) DeleteAdminProductRepository(request IdAdminProductRequest) (entities.Products, error) {
	product := entities.Products{}

	tx := r.db.Delete(&product, "id = ?", request.Id)

	if tx.Error != nil {
		return entities.Products{}, tx.Error
	}

	return product, nil
}
