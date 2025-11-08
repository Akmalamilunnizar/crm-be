package asset

type CreateAdminAssetRequest struct {
	Type         string  `json:"type" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Model        string  `json:"model" validate:"required"`
	SerialNumber string  `json:"serial_number" validate:"required"`
	Date         string  `json:"date" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
	Description  string  `json:"description"`
}

type IdAdminAssetRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminAssetRequest struct {
	IdAdminAssetRequest
	CreateAdminAssetRequest
}
