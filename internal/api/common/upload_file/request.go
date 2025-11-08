package upload_file

import "mime/multipart"

type CreateCommonUploadFileRequest struct {
	Name                   string                `json:"name" validate:"required"`
	Path                   string                `json:"path" validate:"required"`
	File                   *multipart.FileHeader `json:"file" validate:"required"`
	ArchiveInstallationId  string                `json:"archive_installation_id"`
}

type IdCommonUploadFileRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateCommonUploadFileRequest struct {
	IdCommonUploadFileRequest
	CreateCommonUploadFileRequest
}
