package netwatchapi

import (
	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db}
}

// Device operations
func (r *Repo) ListDevices() ([]entities.NetwatchDevice, error) {
	var devices []entities.NetwatchDevice
	err := r.DB.Preload("Customer").Order("created_at DESC").Find(&devices).Error
	return devices, err
}

func (r *Repo) GetDevice(id string) (*entities.NetwatchDevice, error) {
	var device entities.NetwatchDevice
	err := r.DB.Preload("Customer").First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *Repo) CreateDevice(device *entities.NetwatchDevice) (*entities.NetwatchDevice, error) {
	err := r.DB.Create(device).Error
	if err != nil {
		return nil, err
	}
	return r.GetDevice(device.ID)
}

func (r *Repo) UpdateDevice(device *entities.NetwatchDevice) (*entities.NetwatchDevice, error) {
	err := r.DB.Save(device).Error
	if err != nil {
		return nil, err
	}
	return r.GetDevice(device.ID)
}

func (r *Repo) DeleteDevice(id string) error {
	return r.DB.Delete(&entities.NetwatchDevice{}, "id = ?", id).Error
}

// Event operations
func (r *Repo) ListEvents() ([]entities.NetwatchEvent, error) {
	var events []entities.NetwatchEvent
	err := r.DB.Preload("Device").Order("event_time DESC").Find(&events).Error
	return events, err
}

func (r *Repo) GetEvent(id string) (*entities.NetwatchEvent, error) {
	var event entities.NetwatchEvent
	err := r.DB.Preload("Device").First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *Repo) CreateEvent(event *entities.NetwatchEvent) (*entities.NetwatchEvent, error) {
	err := r.DB.Create(event).Error
	if err != nil {
		return nil, err
	}
	return r.GetEvent(event.ID)
}

// Configuration operations
func (r *Repo) ListConfigs() ([]entities.NetwatchConfig, error) {
	var configs []entities.NetwatchConfig
	err := r.DB.Order("created_at DESC").Find(&configs).Error
	return configs, err
}

func (r *Repo) GetConfig(id string) (*entities.NetwatchConfig, error) {
	var config entities.NetwatchConfig
	err := r.DB.First(&config, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *Repo) CreateConfig(config *entities.NetwatchConfig) (*entities.NetwatchConfig, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		return nil, err
	}
	return r.GetConfig(config.ID)
}

func (r *Repo) UpdateConfig(config *entities.NetwatchConfig) (*entities.NetwatchConfig, error) {
	err := r.DB.Save(config).Error
	if err != nil {
		return nil, err
	}
	return r.GetConfig(config.ID)
}

func (r *Repo) DeleteConfig(id string) error {
	return r.DB.Delete(&entities.NetwatchConfig{}, "id = ?", id).Error
}
