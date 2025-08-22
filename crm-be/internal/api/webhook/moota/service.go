package moota

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
)

type WebhookMootaServiceInterface interface {
	GetAllWebhookMootaService() ([]entities.Areas, error)
	CreateWebhookMootaService(request []CreateWebhookMootaRequest) ([]CreateWebhookMootaRequest, error)
	GetByIdWebhookMootaService(request IdWebhookMootaRequest) (entities.Areas, error)
	UpdateWebhookMootaService(request UpdateWebhookMootaRequest) (entities.Areas, error)
	DeleteWebhookMootaService(request IdWebhookMootaRequest) (entities.Areas, error)
}
type WebhookMootaServiceStruct struct {
	repository WebhookMootaRepositoryInterface
}

func NewWebhookMootaService(repository WebhookMootaRepositoryInterface) *WebhookMootaServiceStruct {
	return &WebhookMootaServiceStruct{repository}
}

func (s WebhookMootaServiceStruct) GetAllWebhookMootaService() ([]entities.Areas, error) {
	areas, err := s.repository.FindWebhookMootaRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s WebhookMootaServiceStruct) GetByIdWebhookMootaService(request IdWebhookMootaRequest) (entities.Areas, error) {
	area := entities.Areas{}
	err := copier.Copy(&area, &request)
	if err != nil {
		return area, err
	}

	area, err = s.repository.FindByIdWebhookMootaRepository(request)
	if err != nil {
		return area, err
	}

	return area, nil
}

func (s WebhookMootaServiceStruct) CreateWebhookMootaService(request []CreateWebhookMootaRequest) ([]CreateWebhookMootaRequest, error) {

	return request, nil
}

func (s WebhookMootaServiceStruct) UpdateWebhookMootaService(request UpdateWebhookMootaRequest) (entities.Areas, error) {
	area, err := s.repository.UpdateWebhookMootaRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s WebhookMootaServiceStruct) DeleteWebhookMootaService(request IdWebhookMootaRequest) (entities.Areas, error) {
	area, err := s.repository.DeleteWebhookMootaRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}
