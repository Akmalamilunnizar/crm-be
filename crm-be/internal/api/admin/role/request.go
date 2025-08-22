package role

type CreateAdminRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type IdAdminRoleRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminRoleRequest struct {
	IdAdminRoleRequest
	CreateAdminRoleRequest
}
