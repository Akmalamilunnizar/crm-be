package asset

type CreateAdminAssetRequest struct {
	Type         string  `json:"type" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Model        string  `json:"model" validate:"required"`
	SerialNumber string  `json:"serial_number" validate:"required"`
	Date         string  `json:"date" validate:"required"`
	CompanyID    *string `json:"company_id"`
	Quantity     int64   `json:"quantity" validate:"required"`
	Status       string  `json:"status" `
	Price        int64   `json:"price" validate:"required"`
	Description  string  `json:"description"`
	StatusInOut  string  `json:"status_in_out" validate:"required"`
}

type IdAdminAssetRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminAssetRequest struct {
	IdAdminAssetRequest
	CreateAdminAssetRequest
}
