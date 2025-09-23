package entities

import "time"

// TechnicianTeamMember represents a member of the technician team assigned to a ticket
type TechnicianTeamMember struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID  uint64     `gorm:"column:ticket_id;not null" json:"ticket_id"`
	UserID    string     `gorm:"column:user_id;type:varchar(191);not null" json:"user_id"`
	Role      string     `gorm:"column:role;type:varchar(50);not null" json:"role"` // senior|junior|helper
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (TechnicianTeamMember) TableName() string { return "technician_team_members" }
