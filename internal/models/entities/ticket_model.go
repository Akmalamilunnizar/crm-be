package entities

import (
	"time"

	"gorm.io/gorm"
)

type TroubleType string
type TicketStatus string
type Assignee string

const (
	TroubleWifi     TroubleType = "wifi"
	TroubleInternet TroubleType = "internet"
	TroubleHardware TroubleType = "hardware"
	TroublePower    TroubleType = "power"
	TroubleOther    TroubleType = "other"
	// TroubleNetwork  TroubleType = "network" // New type for Netwatch-triggered tickets - DISABLED FOR TESTING

	// Use role names instead of hardcoded UUIDs - these will be looked up dynamically
	AssignCS   Assignee = "CUSTOMER SERVICE"
	AssignNOC  Assignee = "NOC"
	AssignTech Assignee = "TECHNICIAN"
)

type TroubleTicket struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID      string  `gorm:"type:varchar(191)" json:"customer_id"`
	Type            *string `gorm:"type:varchar(191)" json:"type,omitempty"`
	Title           string  `gorm:"type:longtext" json:"title"`
	Description     *string `gorm:"type:text" json:"description,omitempty"`
	Status          string  `gorm:"type:varchar(191);default:'unfinished'" json:"status"`
	AssignedTo      *string `gorm:"type:varchar(191)" json:"assigned_to,omitempty"`
	CurrentAssignee string  `gorm:"column:current_assignee_role;type:varchar(191)" json:"current_assignee_role"`
	CustomerNote    *string `gorm:"type:text" json:"customer_note,omitempty"`
	NOCNote         *string `gorm:"type:text" json:"noc_note,omitempty"`
	TechnicianNote  *string `gorm:"type:text" json:"technician_note,omitempty"`

	// Image fields
	ImgCS     *string `gorm:"type:varchar(60)" json:"img_cs,omitempty"`
	ImgNOC    *string `gorm:"type:varchar(60)" json:"img_noc,omitempty"`
	ImgTechBF *string `gorm:"type:varchar(60)" json:"img_tech_bf,omitempty"`
	ImgTechAF *string `gorm:"type:varchar(60)" json:"img_tech_af,omitempty"`

	// Netwatch integration fields - DISABLED FOR TESTING
	// CreatedByNetwatch bool            `gorm:"default:false" json:"created_by_netwatch"`
	// NetwatchEventID   *string         `gorm:"type:varchar(36)" json:"netwatch_event_id"`
	// NetwatchEvent     *NetwatchEvent  `json:"netwatch_event,omitempty" gorm:"foreignKey:NetwatchEventID;references:ID"`
	// DeviceID          *string         `gorm:"type:varchar(36)" json:"device_id"`
	// Device            *NetwatchDevice `json:"device,omitempty" gorm:"foreignKey:DeviceID;references:ID"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (TroubleTicket) TableName() string { return "trouble_tickets" }

// Enforce sane defaults before insert regardless of caller
func (t *TroubleTicket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" || t.Status == "received" {
		t.Status = "unfinished"
	}
	// Note: CurrentAssignee will be set by the service layer using dynamic role lookup
	// This hook is kept for other defaults but role assignment is handled in service
	return nil
}
