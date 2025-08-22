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

	AssignCS   Assignee = "1631f7d5-7d01-40af-8d24-1692cefa205a"
	AssignNOC  Assignee = "752e699f-c48c-44ed-8dfa-c8962b4be7ab"
	AssignTech Assignee = "11f0fba5-b49c-4237-9250-ee1b873a7c2b"
)

type TroubleTicket struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID      string     `gorm:"type:varchar(191)" json:"customer_id"`
	Type            *string    `gorm:"type:varchar(191)" json:"type,omitempty"`
	Title           string     `gorm:"type:longtext" json:"title"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	Status          string     `gorm:"type:varchar(191);default:'unfinished'" json:"status"`
	AssignedTo      string     `gorm:"type:varchar(191);not null" json:"assigned_to"`
	CurrentAssignee string     `gorm:"column:current_assignee_role;type:varchar(191)" json:"current_assignee_role"`
	CustomerNote    *string    `gorm:"type:text" json:"customer_note,omitempty"`
	NOCNote         *string    `gorm:"type:text" json:"noc_note,omitempty"`
	TechnicianNote  *string    `gorm:"type:text" json:"technician_note,omitempty"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (TroubleTicket) TableName() string { return "trouble_tickets" }

// Enforce sane defaults before insert regardless of caller
func (t *TroubleTicket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" || t.Status == "received" {
		t.Status = "unfinished"
	}
	if t.CurrentAssignee == "" {
		t.CurrentAssignee = string(AssignCS)
	}
	return nil
}
