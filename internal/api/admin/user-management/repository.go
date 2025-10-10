package usermanagement

import (
	"errors"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"skripsi-be/internal/models/dto"
	"skripsi-be/internal/models/entities"
)

type AdminUserManagementRepositoryInterface interface {
	FindAdminUserManagementRepository(request SearchAdminUserManagementRequest) (*[]dto.UserDTO, error)
	FindByIdAdminUserManagementRepository(request IdAdminUserManagementRequest) (dto.UserDTO, error)
	CreateAdminUserManagementRepository(request CreateAdminUserManagementRequest) (dto.UserDTO, error)
	UpdateAdminUserManagementRepository(request UpdateAdminUserManagementRequest) (dto.UserDTO, error)
	DeleteAdminUserManagementRepository(request IdAdminUserManagementRequest) (dto.UserDTO, error)
	GetUserRolePermissionsRepository(userID string) (map[string]int, error)
}

type AdminUserManagementRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminUserManagementRepository(db *gorm.DB) *AdminUserManagementRepositoryStruct {
	return &AdminUserManagementRepositoryStruct{db}
}
func (r *AdminUserManagementRepositoryStruct) FindAdminUserManagementRepository(request SearchAdminUserManagementRequest) (*[]dto.UserDTO, error) {
	// Simulate a database call
	users := []entities.User{}

	tx := r.db.Preload("Role")
	if request.Role != "" && request.Role != "ALL" {
		tx = tx.Joins("JOIN roles ON roles.id = users.role_id").
			Where("roles.name = ?", request.Role)
	}
	tx = tx.Find(&users)

	// Mapping to dto.User
	userDTOs := &[]dto.UserDTO{}
	copier.Copy(&userDTOs, &users)
	if tx.Error != nil {
		return userDTOs, tx.Error
	}

	return userDTOs, nil
}

func (r *AdminUserManagementRepositoryStruct) FindByIdAdminUserManagementRepository(request IdAdminUserManagementRequest) (dto.UserDTO, error) {
	// Simulate a database call
	user := entities.User{}
	userDto := dto.UserDTO{}
	tx := r.db.Preload("Role").First(&user, "id = ?", request.Id)

	if tx.Error != nil {
		return userDto, tx.Error

	}
	if user.ID == "" {
		return dto.UserDTO{}, gorm.ErrRecordNotFound
	}

	copier.Copy(&userDto, &user)

	return userDto, nil
}

func (r *AdminUserManagementRepositoryStruct) CreateAdminUserManagementRepository(request CreateAdminUserManagementRequest) (dto.UserDTO, error) {
	// Simulate a database call
	user := entities.User{}
	userDto := dto.UserDTO{}
	copier.Copy(&user, &request)

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return userDto, err
	}

	if request.Password != request.PasswordConfirm {
		return userDto, errors.New("Password not match!")

	}

	user.Password = string(password)

	tx := r.db.Create(&user)

	copier.Copy(&userDto, &user)

	if tx.Error != nil {
		return userDto, tx.Error
	}

	return userDto, nil
}

func (r *AdminUserManagementRepositoryStruct) UpdateAdminUserManagementRepository(request UpdateAdminUserManagementRequest) (dto.UserDTO, error) {
	// Simulate a database call
	user := entities.User{}
	userDto := dto.UserDTO{}

	tx := r.db.Preload("Role").First(&user, "id = ?", request.Id)
	if tx.Error != nil {
		return userDto, tx.Error
	}
	copier.Copy(&user, &request)

	// Update password only if provided
	if request.Password != "" {
		if request.Password != request.PasswordConfirm {
			return userDto, errors.New("Password not match!")
		}

		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
		if err != nil {
			return userDto, err
		}

		user.Password = string(password)
	}

	tx = r.db.Save(&user)

	if tx.Error != nil {
		return userDto, tx.Error
	}

	return userDto, nil
}

func (r *AdminUserManagementRepositoryStruct) DeleteAdminUserManagementRepository(request IdAdminUserManagementRequest) (dto.UserDTO, error) {
	// Simulate a database call
	user := entities.User{}
	userDto := dto.UserDTO{}

	tx := r.db.First(&user, "id = ?", request.Id)

	copier.Copy(&userDto, &user)
	if tx.Error != nil {
		return userDto, tx.Error
	}

	tx = r.db.Delete(&user)
	if tx.Error != nil {
		return userDto, tx.Error
	}
	return userDto, nil
}

func (r *AdminUserManagementRepositoryStruct) GetUserRolePermissionsRepository(userID string) (map[string]int, error) {
	// Get user's role first
	user := entities.User{}
	tx := r.db.Preload("Role").First(&user, "id = ?", userID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Get role permissions
	var rolePermissions []entities.RolePermission
	tx = r.db.Where("role_id = ?", user.RoleId).Find(&rolePermissions)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Get features to map feature_id to feature name
	var features []entities.Feature
	tx = r.db.Find(&features)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Create feature name to ID mapping
	featureIDToName := make(map[string]string)
	for _, feature := range features {
		featureIDToName[feature.ID] = feature.Name
	}

	// Create permissions map with feature names as keys
	permissions := make(map[string]int)
	for _, permission := range rolePermissions {
		featureName := featureIDToName[permission.FeatureID]
		if featureName != "" {
			permissions[featureName] = permission.CanAccess
		}
	}

	return permissions, nil
}
