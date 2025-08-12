package account

type CreateAdminAccountRequest struct {
	Name string `json:"name" validate:"required"`
}

type IdAdminAccountRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminAccountRequest struct {
	IdAdminAccountRequest
	CreateAdminAccountRequest
}
