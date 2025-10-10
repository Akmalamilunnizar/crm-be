package feature

import (
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminFeatureRoute(app fiber.Router) {
	db := database.GetDB()
	repository := NewAdminFeatureRepository(db)
	service := NewAdminFeatureService(repository)
	handler := NewAdminFeatureHandler(service)

	app.Use(helpers.VerifyToken)
	app.Get("/", handler.GetAllAdminFeatureHandler)
}

type AdminFeatureHandlerStruct struct {
	service AdminFeatureServiceInterface
}

func NewAdminFeatureHandler(service AdminFeatureServiceInterface) *AdminFeatureHandlerStruct {
	return &AdminFeatureHandlerStruct{service}
}

func (h AdminFeatureHandlerStruct) GetAllAdminFeatureHandler(c *fiber.Ctx) error {
	features, err := h.service.GetAllAdminFeatureService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", features)
}

type AdminFeatureServiceInterface interface {
	GetAllAdminFeatureService() ([]entities.Feature, error)
}

type AdminFeatureServiceStruct struct {
	repository AdminFeatureRepositoryInterface
}

func NewAdminFeatureService(repository AdminFeatureRepositoryInterface) AdminFeatureServiceStruct {
	return AdminFeatureServiceStruct{repository}
}

func (s AdminFeatureServiceStruct) GetAllAdminFeatureService() ([]entities.Feature, error) {
	features, err := s.repository.FindAdminFeatureRepository()
	if err != nil {
		return features, err
	}
	return features, nil
}

type AdminFeatureRepositoryInterface interface {
	FindAdminFeatureRepository() ([]entities.Feature, error)
}

type AdminFeatureRepositoryStruct struct {
	db *gorm.DB
}

func NewAdminFeatureRepository(db *gorm.DB) *AdminFeatureRepositoryStruct {
	return &AdminFeatureRepositoryStruct{db}
}

func (r *AdminFeatureRepositoryStruct) FindAdminFeatureRepository() ([]entities.Feature, error) {
	features := []entities.Feature{}
	tx := r.db.Find(&features)
	if tx.Error != nil {
		return features, tx.Error
	}
	return features, nil
}
