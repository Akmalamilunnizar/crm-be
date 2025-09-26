package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerInstallation model untuk instalasi customer
type CustomerInstallation struct {
	ID                      string     `gorm:"column:id;type:varchar;primaryKey" json:"id"`
	CustomerID              *string    `gorm:"column:customer_id;type:varchar;index:idx_customer_installations_customer_id" json:"customer_id,omitempty"`
	TechnicianID            *string    `gorm:"column:technician_id;type:varchar;index:idx_customer_installations_technician_id" json:"technician_id,omitempty"`
	Status                  string     `gorm:"column:status;type:varchar;default:'pending'" json:"status,omitempty"`
	Notes                   string     `gorm:"column:notes;type:text" json:"notes,omitempty"`
	DocumentType            *string    `gorm:"column:document_type;type:enum('KTP','SIM','Paspor')" json:"document_type,omitempty"`
	DocumentPhoto           *string    `gorm:"column:document_photo;type:varchar(255)" json:"document_photo,omitempty"`
	InstallationType        string     `gorm:"column:installation_type;type:enum('new_installation','maintenance','upgrade','downgrade');default:'new_installation'" json:"installation_type,omitempty"`
	TotalAssetsOut          int        `gorm:"column:total_assets_out;type:int;default:0" json:"total_assets_out,omitempty"`
	TotalAssetsIn           int        `gorm:"column:total_assets_in;type:int;default:0" json:"total_assets_in,omitempty"`
	InstallationCompletedAt *time.Time `gorm:"column:installation_completed_at;type:datetime(3)" json:"installation_completed_at,omitempty"`
	TrialEndDate            *time.Time `gorm:"column:trial_end_date;type:date" json:"trial_end_date,omitempty"`
	ServiceReadyDate        *time.Time `gorm:"column:service_ready_date;type:date" json:"service_ready_date,omitempty"`
	OnAirDate               *time.Time `gorm:"column:on_air_date;type:date" json:"on_air_date,omitempty"`
	CreatedAt               time.Time  `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
	UpdatedAt               time.Time  `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`

	// Relations
	Customer          *Customer          `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer,omitempty"`
	Technician        *User              `gorm:"foreignKey:TechnicianID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"technician,omitempty"`
	Images            []Image            `gorm:"foreignKey:ArchiveInstallationId;references:ID" json:"images,omitempty"`
	AssetTransactions []AssetTransaction `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"asset_transactions,omitempty"`
	NetworkDevices    []NetworkDevice    `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"network_devices,omitempty"`
	CustomerServices  []CustomerService  `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"customer_services,omitempty"`
	Cables            []Cable            `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"cables,omitempty"`
}

func (c *CustomerInstallation) TableName() string {
	return "customer_installations"
}

func (c *CustomerInstallation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
