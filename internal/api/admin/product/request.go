package product

type CreateAdminProductRequest struct {
	Name        string `json:"name" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
	Description string `json:"description"`
}

type IdAdminProductRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminProductRequest struct {
	IdAdminProductRequest
	CreateAdminProductRequest
}
