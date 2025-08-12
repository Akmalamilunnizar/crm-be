package dto

import (
	"skripsi-be/internal/models/entities"
	"time"
)

type CustomerDTO struct {
	ID               string    `json:"id"`
	Address          string    `json:"address"`
	AreaID           string    `json:"area_id"`
	CardIdentition   string    `json:"card_identition"`
	CompanyID        string    `json:"company_id"`
	Email            string    `json:"email"`
	Gender           string    `json:"gender"`
	ProductID        string    `json:"product_id"`
	IPStatic         string    `json:"ip_static"`
	Job              string    `json:"job"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	MacAddress       string    `json:"mac_address"`
	Name             string    `json:"name"`
	NoIdentition     int       `json:"no_identition"`
	Phone            string    `json:"phone"`
	TypeOfService    string    `json:"type_of_service"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	InstallationDate time.Time `json:"installation_date"`
	NextPaymentDate  time.Time `json:"next_payment_date"`
}

type DashboardDTO struct {
	Customer CustomerDTO        `json:"customer"`
	Product  entities.Products  `json:"product"`
	Invoice  []entities.Invoice `json:"invoice"`
}
