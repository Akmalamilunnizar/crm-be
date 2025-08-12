package role

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminRoleRepositoryInterface interface {
	FindAdminRoleRepository() ([]entities.Role, error)
	CreateAdminRoleRepository(request CreateAdminRoleRequest) (entities.Role, error)
	FindByIdAdminRoleRepository(request IdAdminRoleRequest) (entities.Role, error)
	UpdateAdminRoleRepository(request UpdateAdminRoleRequest) (entities.Role, error)
	DeleteAdminRoleRepository(request IdAdminRoleRequest) (entities.Role, error)
}

type AdminRoleRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminRoleRepository(db *gorm.DB) *AdminRoleRepositoryStruct {
	return &AdminRoleRepositoryStruct{db}
}

func (r *AdminRoleRepositoryStruct) FindAdminRoleRepository() ([]entities.Role, error) {

	roles := []entities.Role{}
	tx := r.db.Find(&roles)
	if tx.Error != nil {
		return roles, tx.Error
	}

	return roles, nil

}
func (r *AdminRoleRepositoryStruct) FindByIdAdminRoleRepository(request IdAdminRoleRequest) (entities.Role, error) {
	role := entities.Role{}

	tx := r.db.First(&role, "id = ?", request.Id)
	if tx.Error != nil {
		return role, tx.Error
	}

	return role, nil
}

func (r *AdminRoleRepositoryStruct) CreateAdminRoleRepository(request CreateAdminRoleRequest) (entities.Role, error) {
	role := entities.Role{}
	copier.Copy(&role, &request)
	tx := r.db.Create(&role)
	if tx.Error != nil {
		return role, tx.Error
	}

	return role, nil
}

func (r *AdminRoleRepositoryStruct) UpdateAdminRoleRepository(request UpdateAdminRoleRequest) (entities.Role, error) {
	role := entities.Role{}

	tx := r.db.First(&role, "id = ?", request.Id)
	if tx.Error != nil {
		return role, tx.Error
	}

	role.Name = request.Name

	tx = r.db.Save(&role)
	if tx.Error != nil {
		return role, tx.Error
	}

	return role, nil

}

func (r *AdminRoleRepositoryStruct) DeleteAdminRoleRepository(request IdAdminRoleRequest) (entities.Role, error) {
	role := entities.Role{}

	tx := r.db.First(&role, "id = ?", request.Id)
	if tx.Error != nil {
		return role, tx.Error
	}

	tx = r.db.Delete(&role, "id = ?", request.Id)
	if tx.Error != nil {
		return role, tx.Error
	}

	return role, nil
}
