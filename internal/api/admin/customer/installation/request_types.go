package customerinstallation

// CreateReportInstallationRequest - Request struct untuk multipart form data
type CreateReportInstallationRequest struct {
	// Basic Installation Information
	CustomerID              string `form:"customer_id" binding:"required"`
	TechnicianID            string `form:"technician_id" binding:"required"`
	Status                  string `form:"status"`
	Notes                   string `form:"notes"`
	InstallationType        string `form:"installation_type"`
	OnAirDate               string `form:"on_air_date"`
	TrialEndDate            string `form:"trial_end_date"`
	ServiceReadyDate        string `form:"service_ready_date"`
	InstallationCompletedAt string `form:"installation_completed_at"`

	// Document Information
	DocumentType  string `form:"document_type"`
	DocumentPhoto string `form:"document_photo"` // Will be handled as file upload

	// Network Device Information
	SwitchID             string `form:"switch_id"`
	PortNumber           string `form:"port_number"`
	RemotePort           string `form:"remote_port"`
	EthPort              string `form:"eth_port"`
	MacAddress           string `form:"mac_address"`
	IPStatic             string `form:"ip_static"`
	StatusPerangkat      string `form:"status_perangkat"`
	KepemilikanPerangkat string `form:"kepemilikan_perangkat"`
	LastPingStatus       string `form:"last_ping_status"`

	// Asset Information (via network_devices.assets_id)
	AssetsID string `form:"assets_id" binding:"required"`

	// Cable Information
	CableType   string  `form:"cable_type"`
	CableLength float64 `form:"cable_length"`

	// Customer Service Information
	EndPortType       string `form:"end_port_type"`
	UserLogin         string `form:"user_login"`
	Password          string `form:"password"`
	UserStatus        string `form:"user_status"`
	InstallationNotes string `form:"installation_notes"`
}

// Note: Legacy request structs are defined in the original report files
// This file only contains the new multipart form request struct
