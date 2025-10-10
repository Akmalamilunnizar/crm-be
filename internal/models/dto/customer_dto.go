package dto

import (
	"skripsi-be/internal/models/entities"
	"time"
)

type CustomerDTO struct {
	ID                    string    `json:"id"`
	Address               string    `json:"address"`
	AreaID                string    `json:"area_id"`
	Alias                 string    `json:"alias"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	Name                  string    `json:"name"`
	Phone                 string    `json:"phone"`
	ServiceRequestDate    string    `json:"service_request_date"`
	ProposedPackage       string    `json:"proposed_package"`
	BandwidthCapacity     string    `json:"bandwidth_capacity"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	InstallationDate      time.Time `json:"installation_date"`
	SalesRepresentativeID *string   `json:"sales_representative_id"`
}

type DashboardDTO struct {
	Customer CustomerDTO        `json:"customer"`
	Product  entities.Products  `json:"product"`
	Invoice  []entities.Invoice `json:"invoice"`
}
