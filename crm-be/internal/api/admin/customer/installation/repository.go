package customerinstallation

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"skripsi-be/internal/models/entities"
)

type AdminCustomerInstallationRepositoryInterface interface {
	CreateAdminCustomerInstallationRepository(customer CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	FindByIdAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error)
	FindAdminCustomerInstallationRepository() ([]entities.CustomerInstallation, error)
}
type AdminCustomerInstallationRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminCustomerInstallationRepository(db *gorm.DB) AdminCustomerInstallationRepositoryStruct {
	return AdminCustomerInstallationRepositoryStruct{db}
}

func (r AdminCustomerInstallationRepositoryStruct) FindAdminCustomerInstallationRepository() ([]entities.CustomerInstallation, error) {
	customerInstallations := []entities.CustomerInstallation{}
	tx := r.db.Preload("Customer").Preload("Technician").Preload("Images").Find(&customerInstallations)

	if tx.Error != nil {
		return customerInstallations, tx.Error
	}

	return customerInstallations, nil
}
func (r AdminCustomerInstallationRepositoryStruct) FindByIdAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Customer").Preload("Technician").Preload("Images").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}
func (r AdminCustomerInstallationRepositoryStruct) CreateAdminCustomerInstallationRepository(request CreateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	images := []entities.Image{}
	customerInstallation := entities.CustomerInstallation{}
	copier.Copy(&customerInstallation, &request)

	tx := r.db.Begin()
	txCustomer := tx.Create(&customerInstallation)
	if txCustomer.Error != nil {
		tx.Rollback()
		return entities.CustomerInstallation{}, txCustomer.Error
	}

	txImage := tx.Where("id IN ?", request.ImageIds).Find(&images).Update("archive_installation_id", customerInstallation.ID)
	if txImage.Error != nil {
		tx.Rollback()
		return entities.CustomerInstallation{}, tx.Error
	}

	if tx.Error != nil {
		tx.Rollback()
		return entities.CustomerInstallation{}, tx.Error
	}

	tx.Commit()
	return customerInstallation, nil

}
func (r AdminCustomerInstallationRepositoryStruct) UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Save(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}

func (r AdminCustomerInstallationRepositoryStruct) DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Area").Preload("Company").Preload("Product").Find(&customer, "id = ?", request.Id)
	if tx.Error != nil {
		return customer, tx.Error
	}

	copier.Copy(&customer, &request)

	tx = r.db.Delete(&customer)
	if tx.Error != nil {
		return customer, tx.Error
	}

	return customer, tx.Error
}
