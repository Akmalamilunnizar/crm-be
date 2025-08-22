package upload_file

import (
	"skripsi-be/internal/models/entities"

	"github.com/jinzhu/copier"
)

type CommonUploadFileServiceInterface interface {
	GetAllCommonUploadFileService() ([]entities.Image, error)
	GetByIdCommonUploadFileService(request IdCommonUploadFileRequest) (entities.Image, error)
	CreateCommonUploadFileService(request CreateCommonUploadFileRequest) (entities.Image, error)
	UpdateCommonUploadFileService(request UpdateCommonUploadFileRequest) (entities.Image, error)
	DeleteCommonUploadFileService(request IdCommonUploadFileRequest) (entities.Image, error)
}
type CommonUploadFileServiceStruct struct {
	repository CommonUploadFileRepositoryInterface
}

func NewCommonUploadFileService(repository CommonUploadFileRepositoryInterface) *CommonUploadFileServiceStruct {
	return &CommonUploadFileServiceStruct{repository}
}

func (s CommonUploadFileServiceStruct) GetAllCommonUploadFileService() ([]entities.Image, error) {
	areas, err := s.repository.FindCommonUploadFileRepository()

	if err != nil {
		return areas, err
	}

	return areas, err
}

func (s CommonUploadFileServiceStruct) GetByIdCommonUploadFileService(request IdCommonUploadFileRequest) (entities.Image, error) {
	area := entities.Image{}
	err := copier.Copy(&area, &request)
	if err != nil {
		return area, err
	}

	area, err = s.repository.FindByIdCommonUploadFileRepository(request)
	if err != nil {
		return area, err
	}

	return area, nil
}

func (s CommonUploadFileServiceStruct) CreateCommonUploadFileService(request CreateCommonUploadFileRequest) (entities.Image, error) {
	image, err := s.repository.CreateCommonUploadFileRepository(request)
	if err != nil {
		return image, err
	}

	return image, nil
}

func (s CommonUploadFileServiceStruct) UpdateCommonUploadFileService(request UpdateCommonUploadFileRequest) (entities.Image, error) {
	area, err := s.repository.UpdateCommonUploadFileRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}

func (s CommonUploadFileServiceStruct) DeleteCommonUploadFileService(request IdCommonUploadFileRequest) (entities.Image, error) {
	area, err := s.repository.DeleteCommonUploadFileRepository(request)

	if err != nil {
		return area, err
	}

	return area, nil
}
