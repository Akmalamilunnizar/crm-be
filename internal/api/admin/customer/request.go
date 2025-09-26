package customer

type CreateAdminCustomerRequest struct {
	Name                  string  `json:"name" validate:"required"`
	Alias                 string  `json:"alias"`
	Address               string  `json:"address" validate:"required"`
	AreaID                string  `json:"area_id" validate:"required"`
	Phone                 string  `json:"phone" validate:"required"`
	Latitude              float64 `json:"latitude" validate:"required"`
	Longitude             float64 `json:"longitude" validate:"required"`
	ServiceRequestDate    string  `json:"service_request_date" validate:"required"`
	ProposedPackage       string  `json:"proposed_package" validate:"required"`
	BandwidthCapacity     string  `json:"bandwidth_capacity" validate:"required"`
	SalesRepresentativeID *string `json:"sales_representative_id"`
}

type IdAdminCustomerRequest struct {
	Id string `json:"id"`
}

type UpdateAdminCustomerRequest struct {
	IdAdminCustomerRequest
	CreateAdminCustomerRequest
}
