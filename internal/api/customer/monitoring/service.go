package monitoring

import (
	"fmt"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}


// GetConnectionHistory returns connection history for a customer
func (s *Service) GetConnectionHistory(customerID string, timeRange string) ([]ConnectionEvent, error) {
	// Get the customer's device
	var device entities.NetwatchDevice
	err := s.db.Where("customer_id = ?", customerID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []ConnectionEvent{}, nil
		}
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	// Calculate time range
	var startTime time.Time
	switch timeRange {
	case "1h":
		startTime = time.Now().Add(-1 * time.Hour)
	case "6h":
		startTime = time.Now().Add(-6 * time.Hour)
	case "1d":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// Get events for the device within the time range
	var events []entities.NetwatchEvent
	err = s.db.Where("device_id = ? AND event_time >= ?", device.ID, startTime).
		Order("event_time DESC").
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	// Convert to response format
	var connectionEvents []ConnectionEvent
	for _, event := range events {
		connectionEvents = append(connectionEvents, ConnectionEvent{
			ID:        event.ID,
			Status:    event.EventType,
			Timestamp: event.EventTime,
			Message:   fmt.Sprintf("Connection %s", event.EventType),
		})
	}

	return connectionEvents, nil
}

