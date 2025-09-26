package entities

import "time"

// TicketStep represents a troubleshooting step done by technician
type TicketStep struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID    uint64     `gorm:"column:ticket_id;not null" json:"ticket_id"`
	StepOrder   int        `gorm:"column:step_order;not null" json:"step_order"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	ImagePath   *string    `gorm:"column:image_path;type:varchar(191)" json:"image_path,omitempty"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (TicketStep) TableName() string { return "ticket_steps" }
