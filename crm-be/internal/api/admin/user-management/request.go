package usermanagement

type CreateAdminUserManagementRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	LogoUrl         string `json:"logo_url"`
	RoleId          string `json:"role_id" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

type IdAdminUserManagementRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminUserManagementRequest struct {
	IdAdminUserManagementRequest
	CreateAdminUserManagementRequest
}

type SearchAdminUserManagementRequest struct {
	Role string `json:"role"`
}
