package account

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AdminAccountRepositoryInterface interface {
	FindAdminAccountRepository() ([]entities.Accounts, error)
	CreateAdminAccountRepository(request CreateAdminAccountRequest) (entities.Accounts, error)
	FindByIdAdminAccountRepository(request IdAdminAccountRequest) (entities.Accounts, error)
	UpdateAdminAccountRepository(request UpdateAdminAccountRequest) (entities.Accounts, error)
	DeleteAdminAccountRepository(request IdAdminAccountRequest) (entities.Accounts, error)
}

type AdminAccountRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminAccountRepository(db *gorm.DB) *AdminAccountRepositoryStruct {
	return &AdminAccountRepositoryStruct{db}
}

func (r AdminAccountRepositoryStruct) FindAdminAccountRepository() ([]entities.Accounts, error) {
	accounts := []entities.Accounts{}
	tx := r.db.Find(&accounts)

	if tx.Error != nil {
		return accounts, tx.Error
	}

	return accounts, nil
}

func (r AdminAccountRepositoryStruct) FindByIdAdminAccountRepository(request IdAdminAccountRequest) (entities.Accounts, error) {
	accounts := entities.Accounts{}
	tx := r.db.First(&accounts, "id = ?", request.Id)

	if tx.Error != nil {
		return accounts, tx.Error
	}

	return accounts, nil
}

func (r AdminAccountRepositoryStruct) CreateAdminAccountRepository(request CreateAdminAccountRequest) (entities.Accounts, error) {
	account := entities.Accounts{}
	copier.Copy(&account, &request)
	tx := r.db.Create(&account)
	if tx.Error != nil {
		return entities.Accounts{}, tx.Error
	}

	return account, nil
}

func (r AdminAccountRepositoryStruct) UpdateAdminAccountRepository(request UpdateAdminAccountRequest) (entities.Accounts, error) {
	accounts := entities.Accounts{}
	tx := r.db.First(&accounts, "id = ?", request.Id)

	if tx.Error != nil {
		return accounts, tx.Error
	}
	copier.Copy(&accounts, &request)

	r.db.Save(&accounts)
	return accounts, nil
}

func (r AdminAccountRepositoryStruct) DeleteAdminAccountRepository(request IdAdminAccountRequest) (entities.Accounts, error) {

	accounts := entities.Accounts{}

	tx := r.db.First(&accounts, "id = ?", request.Id)
	if tx.Error != nil {
		return accounts, tx.Error
	}

	tx = r.db.Delete(&accounts)
	if tx.Error != nil {
		return accounts, tx.Error
	}
	return accounts, nil
}
