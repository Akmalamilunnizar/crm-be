package networkdevice

type CreateNetworkDeviceRequest struct {
	CustomerID      string  `json:"customer_id" validate:"required"`
	IPStatic        *string `json:"ip_static"`
	MacAddress      *string `json:"mac_address"`
	StatusPerangkat string  `json:"status_perangkat" validate:"required"`
	LastPingStatus  string  `json:"last_ping_status" validate:"required"`
	AssetsID        *string `json:"assets_id"`
	ProductID       *string `json:"product_id"`
}

type UpdateNetworkDeviceRequest struct {
	ID              string  `json:"id" validate:"required"`
	CustomerID      string  `json:"customer_id" validate:"required"`
	IPStatic        *string `json:"ip_static"`
	MacAddress      *string `json:"mac_address"`
	StatusPerangkat string  `json:"status_perangkat" validate:"required"`
	LastPingStatus  string  `json:"last_ping_status" validate:"required"`
	AssetsID        *string `json:"assets_id"`
	ProductID       *string `json:"product_id"`
}

type IdNetworkDeviceRequest struct {
	Id string `json:"id"`
}
