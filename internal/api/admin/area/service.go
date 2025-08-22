package area

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
)

type AdminAreaServiceInterface interface {
	GetAllAdminAreaService() ([]entities.Areas, error)
	CreateAdminAreaService(request CreateAdminAreaRequest) (entities.Areas, error)
	GetByIdAdminAreaService(request IdAdminAreaRequest) (entities.Areas, error)
	UpdateAdminAreaService(request UpdateAdminAreaRequest) (entities.Areas, error)
	DeleteAdminAreaService(request IdAdminAreaRequest) (entities.Areas, error)
}
type AdminAreaServiceStruct struct {
	repository AdminAreaRepositoryInterface
}

func NewAdminAreaService(repository AdminAreaRepositoryInterface) *AdminAreaServiceStruct {
	return &AdminAreaServiceStruct{repository}
}

func (s AdminAreaServiceStruct) GetAllAdminAreaService() ([]entities.Areas, error) {
	areas, err := s.repository.FindAdminAreaRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s AdminAreaServiceStruct) GetByIdAdminAreaService(request IdAdminAreaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	err := copier.Copy(&area, &request)
	if err != nil {
		return area, err
	}

	area, err = s.repository.FindByIdAdminAreaRepository(request)
	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAreaServiceStruct) CreateAdminAreaService(request CreateAdminAreaRequest) (entities.Areas, error) {

	area, err := s.repository.CreateAdminAreaRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAreaServiceStruct) UpdateAdminAreaService(request UpdateAdminAreaRequest) (entities.Areas, error) {
	area, err := s.repository.UpdateAdminAreaRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s AdminAreaServiceStruct) DeleteAdminAreaService(request IdAdminAreaRequest) (entities.Areas, error) {
	area, err := s.repository.DeleteAdminAreaRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}
