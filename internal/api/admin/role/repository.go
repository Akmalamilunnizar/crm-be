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

	tx := r.db.Preload("RolePermissions").First(&role, "id = ?", request.Id)
	if tx.Error != nil {
		return role, tx.Error
	}

	return role, nil
}

func (r *AdminRoleRepositoryStruct) CreateAdminRoleRepository(request CreateAdminRoleRequest) (entities.Role, error) {
	role := entities.Role{}
	copier.Copy(&role, &request)
	
	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return role, tx.Error
	}
	
	// Create role
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return role, err
	}
	
	// Create role permissions
	for _, permission := range request.RolePermissions {
		permission.RoleID = role.ID
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			return role, err
		}
	}
	
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return role, err
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

	// Start transaction for update
	dbTx := r.db.Begin()
	if dbTx.Error != nil {
		return role, dbTx.Error
	}

	// Update role
	if err := dbTx.Save(&role).Error; err != nil {
		dbTx.Rollback()
		return role, err
	}

	// Delete existing role permissions
	if err := dbTx.Where("role_id = ?", role.ID).Delete(&entities.RolePermission{}).Error; err != nil {
		dbTx.Rollback()
		return role, err
	}

	// Create new role permissions
	for _, permission := range request.RolePermissions {
		permission.RoleID = role.ID
		if err := dbTx.Create(&permission).Error; err != nil {
			dbTx.Rollback()
			return role, err
		}
	}

	// Commit transaction
	if err := dbTx.Commit().Error; err != nil {
		return role, err
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
