package role

import "skripsi-be/internal/models/entities"

type CreateAdminRoleRequest struct {
	Name            string                        `json:"name" validate:"required"`
	RolePermissions []entities.RolePermission     `json:"role_permissions"`
}

type IdAdminRoleRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminRoleRequest struct {
	IdAdminRoleRequest
	CreateAdminRoleRequest
}
