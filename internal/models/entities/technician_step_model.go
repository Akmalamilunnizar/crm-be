package entities

import "time"

// TechnicianStep represents a predefined step in the technician workflow
type TechnicianStep struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StepOrder   int        `gorm:"column:step_order;not null" json:"step_order"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Tools       string     `gorm:"type:text" json:"tools"`       // Peralatan (Tools)
	SpareParts  string     `gorm:"type:text" json:"spare_parts"` // Suku Cadang (Spare Parts)
	Procedure   string     `gorm:"type:text" json:"procedure"`   // Langkah Pengecekan (Procedure)
	Solution    string     `gorm:"type:text" json:"solution"`    // Solusi (Solution)
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (TechnicianStep) TableName() string { return "technician_steps" }

// TicketTechnicianStep represents a technician's progress on a specific step for a ticket
type TicketTechnicianStep struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID       uint64     `gorm:"column:ticket_id;not null" json:"ticket_id"`
	StepID         uint64     `gorm:"column:step_id;not null" json:"step_id"`
	TechnicianID   string     `gorm:"column:technician_id;type:varchar(191);not null" json:"technician_id"`
	Status         string     `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, done, not_applicable, needs_spare_parts
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`
	SparePartsUsed *string    `gorm:"type:text" json:"spare_parts_used,omitempty"` // JSON array of spare parts
	ImagePaths     *string    `gorm:"type:text" json:"image_paths,omitempty"`      // JSON array of image paths
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`

	// Relations
	Step       *TechnicianStep `json:"step,omitempty" gorm:"foreignKey:StepID;references:ID"`
	Technician *User           `json:"technician,omitempty" gorm:"foreignKey:TechnicianID;references:ID"`
}

func (TicketTechnicianStep) TableName() string { return "ticket_technician_steps" }

// SparePart represents available spare parts
type SparePart struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Category    string     `gorm:"type:varchar(100)" json:"category"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (SparePart) TableName() string { return "spare_parts" }
