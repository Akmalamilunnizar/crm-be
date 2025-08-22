package moota

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type WebhookMootaRepositoryInterface interface {
	FindWebhookMootaRepository() ([]entities.Areas, error)
	CreateWebhookMootaRepository(request CreateWebhookMootaRequest) (entities.Areas, error)
	FindByIdWebhookMootaRepository(request IdWebhookMootaRequest) (entities.Areas, error)
	UpdateWebhookMootaRepository(request UpdateWebhookMootaRequest) (entities.Areas, error)
	DeleteWebhookMootaRepository(request IdWebhookMootaRequest) (entities.Areas, error)
}

type WebhookMootaRepositoryStruct struct {
	db *gorm.DB
}

func NewWebhookMootaRepository(db *gorm.DB) *WebhookMootaRepositoryStruct {
	return &WebhookMootaRepositoryStruct{db}
}

func (r WebhookMootaRepositoryStruct) FindWebhookMootaRepository() ([]entities.Areas, error) {
	areas := []entities.Areas{}
	tx := r.db.Find(&areas)

	if tx.Error != nil {
		return areas, tx.Error
	}

	return areas, nil
}

func (r WebhookMootaRepositoryStruct) FindByIdWebhookMootaRepository(request IdWebhookMootaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	tx := r.db.First(&area, "id = ?", request.Id)

	if tx.Error != nil {
		return area, tx.Error
	}

	return area, nil
}

func (r WebhookMootaRepositoryStruct) CreateWebhookMootaRepository(request CreateWebhookMootaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	copier.Copy(&area, &request)
	tx := r.db.Create(&area)
	if tx.Error != nil {
		return entities.Areas{}, tx.Error
	}

	return area, nil
}

func (r WebhookMootaRepositoryStruct) UpdateWebhookMootaRepository(request UpdateWebhookMootaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	tx := r.db.First(&area, "id = ?", request.Id)

	if tx.Error != nil {
		return area, tx.Error
	}
	copier.Copy(&area, &request)

	r.db.Save(&area)
	return area, nil
}

func (r WebhookMootaRepositoryStruct) DeleteWebhookMootaRepository(request IdWebhookMootaRequest) (entities.Areas, error) {

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
