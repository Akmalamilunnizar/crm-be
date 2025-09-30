package dto

import (
	"time"
)

type ReportInternet struct {
	Name             string    `json:"name"`
	InstallationDate time.Time `json:"installation_date"`
	ProductName      string    `json:"product"`
}
