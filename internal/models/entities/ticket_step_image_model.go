package entities

import "time"

type TicketStepImage struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StepID    uint64     `gorm:"column:step_id;not null" json:"step_id"`
	Path      string     `gorm:"column:path;type:varchar(191);not null" json:"path"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

func (TicketStepImage) TableName() string { return "ticket_step_images" }
