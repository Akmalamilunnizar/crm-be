package product

import "skripsi-be/internal/models/entities"

type AdminProductServiceInterface interface {
}

type AdminProductServiceStruct struct {
	repository AdminProductRepositoryStruct
}

func NewAdminProductServiceStruct(repository AdminProductRepositoryStruct) AdminProductServiceStruct {
	return AdminProductServiceStruct{repository}
}

func (s AdminProductServiceStruct) GetAllAdminProductService() ([]entities.Products, error) {
	products, err := s.repository.FindAdminProductRepository()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s AdminProductServiceStruct) GetByIdAdminProductService(id IdAdminProductRequest) (entities.Products, error) {
	product, err := s.repository.FindByIdAdminProductRepository(id)
	if err != nil {
		return entities.Products{}, err
	}

	return product, nil
}

func (s AdminProductServiceStruct) CreateAdminProductService(request CreateAdminProductRequest) (entities.Products, error) {
	product, err := s.repository.CreateAdminProductRepository(request)
	if err != nil {
		return entities.Products{}, err
	}

	return product, nil
}

func (s AdminProductServiceStruct) UpdateAdminProductService(request UpdateAdminProductRequest) (entities.Products, error) {
	product, err := s.repository.UpdateAdminProductRepository(request)
	if err != nil {
		return entities.Products{}, err

	}

	return product, nil
}
func (s AdminProductServiceStruct) DeleteAdminProductService(request IdAdminProductRequest) (entities.Products, error) {
	product, err := s.repository.DeleteAdminProductRepository(request)
	if err != nil {
		return entities.Products{}, err

	}

	return product, nil
}
