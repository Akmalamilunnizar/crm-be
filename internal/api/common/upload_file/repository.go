package upload_file

import (
	"errors"
	"path/filepath"
	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type CommonUploadFileRepositoryInterface interface {
	FindCommonUploadFileRepository() ([]entities.Image, error)
	CreateCommonUploadFileRepository(request CreateCommonUploadFileRequest) (entities.Image, error)
	FindByIdCommonUploadFileRepository(request IdCommonUploadFileRequest) (entities.Image, error)
	UpdateCommonUploadFileRepository(request UpdateCommonUploadFileRequest) (entities.Image, error)
	DeleteCommonUploadFileRepository(request IdCommonUploadFileRequest) (entities.Image, error)
}

type CommonUploadFileRepositoryStruct struct {
	db *gorm.DB
}

func NewCommonUploadFileRepository(db *gorm.DB) *CommonUploadFileRepositoryStruct {
	return &CommonUploadFileRepositoryStruct{db}
}

func (r CommonUploadFileRepositoryStruct) FindCommonUploadFileRepository() ([]entities.Image, error) {
	areas := []entities.Image{}
	tx := r.db.Find(&areas)

	if tx.Error != nil {
		return areas, tx.Error
	}

	return areas, nil
}

func (r CommonUploadFileRepositoryStruct) FindByIdCommonUploadFileRepository(request IdCommonUploadFileRequest) (entities.Image, error) {
	area := entities.Image{}

	return area, nil
}

func (r CommonUploadFileRepositoryStruct) CreateCommonUploadFileRepository(request CreateCommonUploadFileRequest) (entities.Image, error) {

	image := entities.Image{}
	extension := filepath.Ext(request.File.Filename)
	if extension == "" {
		return image, errors.New("File has no extension")
	}
	image.File = request.Name + extension
	image.FullPath = "uploads/" + request.Path + "/" + request.Name + extension
	tx := r.db.Create(&image)
	if tx.Error != nil {
		return image, tx.Error
	}

	return image, nil
}

func (r CommonUploadFileRepositoryStruct) UpdateCommonUploadFileRepository(request UpdateCommonUploadFileRequest) (entities.Image, error) {
	area := entities.Image{}

	return area, nil
}

func (r CommonUploadFileRepositoryStruct) DeleteCommonUploadFileRepository(request IdCommonUploadFileRequest) (entities.Image, error) {

	image := entities.Image{}

	tx := r.db.First(&image, "id = ?", request.Id)
	if tx.Error != nil {
		return image, tx.Error
	}

	tx = r.db.Delete(&image)
	if tx.Error != nil {
		return image, tx.Error
	}
	return image, nil
}
