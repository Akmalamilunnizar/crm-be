package customerinstallation

type CreateAdminCustomerInstallationRequest struct {
	CustomerId   string   `json:"customer_id" validate:"required"`
	TechnicianId string   `json:"technician_id" validate:"required"`
	Date         string   `json:"date" validate:"required"`
	Description  string   `json:"description"`
	ImageIds     []string `json:"image_ids" validate:"required"`
}

type IdAdminCustomerInstallationRequest struct {
	Id string `json:"id"`
}

type UpdateAdminCustomerInstallationRequest struct {
	IdAdminCustomerInstallationRequest
	CreateAdminCustomerInstallationRequest
}
