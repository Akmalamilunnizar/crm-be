package company

type CreateAdminCompanyRequest struct {
	Name        string `json:"name" validate:"required"`
	URL         string `json:"url"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	LogoUrl     string `json:"logo_url"`
	Description string `json:"description"`
	Npwp        string `json:"npwp" validate:"required"`
	Address     string `json:"address" validate:"required"`
}

type IdAdminCompanyRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminCompanyRequest struct {
	IdAdminCompanyRequest
	CreateAdminCompanyRequest
}
