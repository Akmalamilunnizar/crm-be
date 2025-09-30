package customerinstallation

type CreateAdminCustomerInstallationRequest struct {
	CustomerId    string   `json:"customer_id" validate:"required"`
	TechnicianId  string   `json:"technician_id" validate:"required"`
	Status        string   `json:"status"`
	Notes         string   `json:"notes"`
	DocumentType  string   `json:"document_type" validate:"omitempty,oneof=KTP SIM Paspor"`
	DocumentPhoto string   `json:"document_photo"`
	ImageIds      []string `json:"image_ids" validate:"required"`
	OnAirDate     string   `json:"on_air_date"`
}

type IdAdminCustomerInstallationRequest struct {
	Id string `json:"id"`
}

type UpdateAdminCustomerInstallationRequest struct {
	IdAdminCustomerInstallationRequest
	CreateAdminCustomerInstallationRequest
}
