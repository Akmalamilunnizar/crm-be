package area

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminAreaRepositoryInterface interface {
	FindAdminAreaRepository() ([]entities.Areas, error)
	CreateAdminAreaRepository(request CreateAdminAreaRequest) (entities.Areas, error)
	FindByIdAdminAreaRepository(request IdAdminAreaRequest) (entities.Areas, error)
	UpdateAdminAreaRepository(request UpdateAdminAreaRequest) (entities.Areas, error)
	DeleteAdminAreaRepository(request IdAdminAreaRequest) (entities.Areas, error)
}

type AdminAreaRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminAreaRepository(db *gorm.DB) *AdminAreaRepositoryStruct {
	return &AdminAreaRepositoryStruct{db}
}

func (r AdminAreaRepositoryStruct) FindAdminAreaRepository() ([]entities.Areas, error) {
	areas := []entities.Areas{}
	tx := r.db.Find(&areas)

	if tx.Error != nil {
		return areas, tx.Error
	}

	return areas, nil
}

func (r AdminAreaRepositoryStruct) FindByIdAdminAreaRepository(request IdAdminAreaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	tx := r.db.First(&area, "id = ?", request.Id)

	if tx.Error != nil {
		return area, tx.Error
	}

	return area, nil
}

func (r AdminAreaRepositoryStruct) CreateAdminAreaRepository(request CreateAdminAreaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	copier.Copy(&area, &request)
	tx := r.db.Create(&area)
	if tx.Error != nil {
		return entities.Areas{}, tx.Error
	}

	return area, nil
}

func (r AdminAreaRepositoryStruct) UpdateAdminAreaRepository(request UpdateAdminAreaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	tx := r.db.First(&area, "id = ?", request.Id)

	if tx.Error != nil {
		return area, tx.Error
	}
	copier.Copy(&area, &request)

	r.db.Save(&area)
	return area, nil
}

func (r AdminAreaRepositoryStruct) DeleteAdminAreaRepository(request IdAdminAreaRequest) (entities.Areas, error) {

	area := entities.Areas{}

	tx := r.db.First(&area, "id = ?", request.Id)
	if tx.Error != nil {
		return area, tx.Error
	}

	tx = r.db.Delete(&area)
	if tx.Error != nil {
		return area, tx.Error
	}
	return area, nil
}
