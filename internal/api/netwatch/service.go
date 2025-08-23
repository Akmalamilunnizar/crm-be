package netwatchapi

import (
	"skripsi-be/internal/models/entities"
	"skripsi-be/internal/services"
)

type Service struct {
	repo            *Repo
	netwatchService *services.NetwatchService
}

func NewService(r *Repo, netwatchService *services.NetwatchService) *Service {
	return &Service{
		repo:            r,
		netwatchService: netwatchService,
	}
}

// Device management
func (s *Service) ListDevices() ([]entities.NetwatchDevice, error) {
	return s.repo.ListDevices()
}

func (s *Service) GetDevice(id string) (*entities.NetwatchDevice, error) {
	return s.repo.GetDevice(id)
}

func (s *Service) CreateDevice(device *entities.NetwatchDevice) (*entities.NetwatchDevice, error) {
	return s.repo.CreateDevice(device)
}

func (s *Service) UpdateDevice(device *entities.NetwatchDevice) (*entities.NetwatchDevice, error) {
	return s.repo.UpdateDevice(device)
}

func (s *Service) DeleteDevice(id string) error {
	return s.repo.DeleteDevice(id)
}

// Event management
func (s *Service) ListEvents() ([]entities.NetwatchEvent, error) {
	return s.repo.ListEvents()
}

func (s *Service) GetEvent(id string) (*entities.NetwatchEvent, error) {
	return s.repo.GetEvent(id)
}

func (s *Service) CreateEvent(event *entities.NetwatchEvent) (*entities.NetwatchEvent, error) {
	createdEvent, err := s.repo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	// Process the event immediately
	if err := s.netwatchService.ProcessNetwatchEvent(createdEvent); err != nil {
		// Log error but don't fail the request
		// The event will be processed later by the background service
	}

	return createdEvent, nil
}

// Configuration management
func (s *Service) ListConfigs() ([]entities.NetwatchConfig, error) {
	return s.repo.ListConfigs()
}

func (s *Service) GetConfig(id string) (*entities.NetwatchConfig, error) {
	return s.repo.GetConfig(id)
}

func (s *Service) CreateConfig(config *entities.NetwatchConfig) (*entities.NetwatchConfig, error) {
	return s.repo.CreateConfig(config)
}

func (s *Service) UpdateConfig(config *entities.NetwatchConfig) (*entities.NetwatchConfig, error) {
	return s.repo.UpdateConfig(config)
}

func (s *Service) DeleteConfig(id string) error {
	return s.repo.DeleteConfig(id)
}

// Sync devices from MikroTik
func (s *Service) SyncDevicesFromMikroTik(configID string) error {
	config, err := s.repo.GetConfig(configID)
	if err != nil {
		return err
	}

	return s.netwatchService.SyncDevicesFromMikroTik(config)
}

// Get monitoring status
func (s *Service) GetMonitoringStatus() (map[string]interface{}, error) {
	devices, err := s.repo.ListDevices()
	if err != nil {
		return nil, err
	}

	events, err := s.repo.ListEvents()
	if err != nil {
		return nil, err
	}

	// Count devices by status
	upCount := 0
	downCount := 0
	for _, device := range devices {
		if device.Status == "up" {
			upCount++
		} else {
			downCount++
		}
	}

	// Count unprocessed events
	unprocessedCount := 0
	for _, event := range events {
		if !event.Processed {
			unprocessedCount++
		}
	}

	return map[string]interface{}{
		"total_devices":      len(devices),
		"up_devices":         upCount,
		"down_devices":       downCount,
		"total_events":       len(events),
		"unprocessed_events": unprocessedCount,
		"monitoring_active":  true, // TODO: Get from actual monitoring status
	}, nil
}
