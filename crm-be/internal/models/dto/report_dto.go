package dto

import (
	"encoding/json"
	"skripsi-be/internal/models/entities"
	"time"
)

type ReportInternet struct {
	Name             string            `json:"name"`
	InstallationDate time.Time         `json:"installation_date"`
	Product          entities.Products `json:"product"`
}

func (r ReportInternet) MarshalJSON() ([]byte, error) {
	type Alias ReportInternet
	return json.Marshal(struct {
		Alias
		Product string `json:"product"`
	}{
		Alias:   Alias(r),
		Product: r.Product.Name,
	})
}
