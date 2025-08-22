package customer

type CreateAdminCustomerRequest struct {
	TypeOfService  string  `json:"type_of_service" validate:"required"`
	Email          string  `json:"email" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	CompanyId      string  `json:"company_id" validate:"required"`
	Gender         string  `json:"gender" validate:"required"`
	CardIdentition string  `json:"card_identition" validate:"required"`
	NoIdentition   int64   `json:"no_identition" validate:"required"`
	AreaID         string  `json:"area_id" validate:"required"`
	Phone          string  `json:"phone" validate:"required"`
	Address        string  `json:"address" validate:"required"`
	Latitude       float64 `json:"latitude" validate:"required"`
	Longitude      float64 `json:"longitude" validate:"required"`
	ProductID      string  `json:"product_id" validate:"required"`
	IpStatic       string  `json:"ip_static" validate:"required"`
	MacAddress     string  `json:"mac_address" validate:"required"`
	Job            string  `json:"job"`
}

type IdAdminCustomerRequest struct {
	Id string `json:"id"`
}

type UpdateAdminCustomerRequest struct {
	IdAdminCustomerRequest
	CreateAdminCustomerRequest
}
