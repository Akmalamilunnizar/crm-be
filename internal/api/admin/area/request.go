package area

type CreateAdminAreaRequest struct {
	NameCity        string `json:"name_city" validate:"required"`
	NameSubdistrict string `json:"name_subdistrict" validate:"required"`
	NameVillage     string `json:"name_village" validate:"required"`
}

type IdAdminAreaRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminAreaRequest struct {
	IdAdminAreaRequest
	CreateAdminAreaRequest
}
