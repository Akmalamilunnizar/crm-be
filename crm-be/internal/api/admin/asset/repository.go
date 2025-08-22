package asset

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminAssetRepositoryInterface interface {
	FindAdminAssetRepository() ([]entities.Asset, error)
	CreateAdminAssetRepository(request CreateAdminAssetRequest) (entities.Asset, error)
	FindByIdAdminAssetRepository(request IdAdminAssetRequest) (entities.Asset, error)
	UpdateAdminAssetRepository(request UpdateAdminAssetRequest) (entities.Asset, error)
	DeleteAdminAssetRepository(request IdAdminAssetRequest) (entities.Asset, error)
}

type AdminAssetRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminAssetRepository(db *gorm.DB) *AdminAssetRepositoryStruct {
	return &AdminAssetRepositoryStruct{db}
}

func (r AdminAssetRepositoryStruct) FindAdminAssetRepository() ([]entities.Asset, error) {
	assets := []entities.Asset{}
	tx := r.db.Find(&assets)
	if tx.Error != nil {
		return assets, nil
	}

	return assets, nil
}
func (r AdminAssetRepositoryStruct) CreateAdminAssetRepository(request CreateAdminAssetRequest) (entities.Asset, error) {
	asset := entities.Asset{}
	copier.Copy(&asset, &request)

	tx := r.db.Create(&asset)
	if tx.Error != nil {
		return asset, tx.Error
	}

	return asset, nil
}
func (r AdminAssetRepositoryStruct) FindByIdAdminAssetRepository(request IdAdminAssetRequest) (entities.Asset, error) {
	asset := entities.Asset{}
	tx := r.db.First(&asset, "id = ?", request.Id)
	if tx.Error != nil {
		return asset, tx.Error
	}
	return asset, nil
}
func (r AdminAssetRepositoryStruct) UpdateAdminAssetRepository(request UpdateAdminAssetRequest) (entities.Asset, error) {
	asset := entities.Asset{}
	tx := r.db.First(&asset, "id = ?", request.Id)
	if tx.Error != nil {
		return asset, tx.Error
	}
	copier.Copy(&asset, &request.CreateAdminAssetRequest)

	tx = r.db.Save(&asset)
	if tx.Error != nil {
		return asset, tx.Error
	}

	return asset, nil
}
func (r AdminAssetRepositoryStruct) DeleteAdminAssetRepository(request IdAdminAssetRequest) (entities.Asset, error) {
	asset := entities.Asset{}
	copier.Copy(&asset, &request)
	tx := r.db.First(&asset, "id = ?", request.Id)
	if tx.Error != nil {
		return asset, tx.Error
	}

	tx = r.db.Delete(&asset)
	if tx.Error != nil {
		return asset, tx.Error
	}

	return asset, nil
}
