package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NetworkDevice represents a network device assigned to customers
type NetworkDevice struct {
	ID                     string                `json:"id" gorm:"primaryKey;type:varchar(191)"`
	CustomerID             string                `json:"customer_id" gorm:"type:varchar(191);not null"`
	Customer               *Customer             `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:ID"`
	AssetsID               string                `json:"assets_id" gorm:"type:varchar(191);not null"`
	Assets                 *Asset                `json:"assets,omitempty" gorm:"foreignKey:AssetsID;references:ID"`
	CustomerInstallationID *string               `json:"customer_installation_id,omitempty" gorm:"type:varchar(191);index:idx_network_devices_customer_installation_id"`
	CustomerInstallation   *CustomerInstallation `json:"customer_installation,omitempty" gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	SwitchID               *string               `json:"switch_id" gorm:"type:varchar(191)"`
	PortNumber             *string               `json:"port_number" gorm:"type:varchar(50)"`
	RemotePort             *string               `json:"remote_port" gorm:"type:varchar(50)"`
	EthPort                *string               `json:"eth_port" gorm:"type:varchar(50)"`
	KepemilikanPerangkat   string                `json:"kepemilikan_perangkat" gorm:"type:enum('owned','leased','customer');default:'owned'"`
	StatusPerangkat        string                `json:"status_perangkat" gorm:"type:enum('active','inactive','maintenance','faulty');default:'active'"`
	LastPingStatus         string                `json:"last_ping_status" gorm:"type:enum('up','down','unknown');default:'unknown'"`
	LastPingTimestamp      *time.Time            `json:"last_ping_timestamp"`
	CreatedAt              time.Time             `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt              time.Time             `json:"updated_at" gorm:"default:current_timestamp"`
	MacAddress             *string               `json:"mac_address" gorm:"type:varchar(191)"`
	IPStatic               *string               `json:"ip_static" gorm:"type:varchar(191)"`
	ProductID              *string               `json:"product_id" gorm:"type:varchar(191)"`
	Product                *Products             `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

func (n *NetworkDevice) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

func (NetworkDevice) TableName() string {
	return "network_devices"
}

