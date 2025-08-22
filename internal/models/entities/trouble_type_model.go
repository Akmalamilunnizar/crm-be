package entities

type TroubleTypeRow struct {
	ID   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

func (TroubleTypeRow) TableName() string { return "trouble_type" }
