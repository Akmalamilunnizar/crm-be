package asset

import "skripsi-be/internal/models/entities"

type AdminAssetServiceInterface interface {
	GetAllAdminAssetService() ([]entities.Asset, error)
	CreateAdminAssetService(request CreateAdminAssetRequest) (entities.Asset, error)
	GetByIdAdminAssetService(request IdAdminAssetRequest) (entities.Asset, error)
	UpdateAdminAssetService(request UpdateAdminAssetRequest) (entities.Asset, error)
	DeleteAdminAssetService(request IdAdminAssetRequest) (entities.Asset, error)
}
type AdminAssetServiceStruct struct {
	repository AdminAssetRepositoryInterface
}

func NewAdminAssetService(repository AdminAssetRepositoryInterface) AdminAssetServiceStruct {
	return AdminAssetServiceStruct{repository}
}

func (s AdminAssetServiceStruct) GetAllAdminAssetService() ([]entities.Asset, error) {
	assets, err := s.repository.FindAdminAssetRepository()
	if err != nil {
		return assets, err
	}
	return assets, err
}
func (s AdminAssetServiceStruct) GetByIdAdminAssetService(request IdAdminAssetRequest) (entities.Asset, error) {
	asset, err := s.repository.FindByIdAdminAssetRepository(request)
	if err != nil {
		return asset, err
	}
	return asset, err
}
func (s AdminAssetServiceStruct) CreateAdminAssetService(request CreateAdminAssetRequest) (entities.Asset, error) {
	asset, err := s.repository.CreateAdminAssetRepository(request)
	if err != nil {
		return asset, err
	}
	return asset, err
}

func (s AdminAssetServiceStruct) UpdateAdminAssetService(request UpdateAdminAssetRequest) (entities.Asset, error) {
	asset, err := s.repository.UpdateAdminAssetRepository(request)
	if err != nil {
		return asset, err
	}
	return asset, err
}
func (s AdminAssetServiceStruct) DeleteAdminAssetService(request IdAdminAssetRequest) (entities.Asset, error) {
	asset, err := s.repository.DeleteAdminAssetRepository(request)
	if err != nil {
		return asset, err
	}
	return asset, err
}
