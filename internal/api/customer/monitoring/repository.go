package monitoring

import (
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// GetCustomerDevice returns the netwatch device for a customer
func (r *Repo) GetCustomerDevice(customerID string) (*entities.NetwatchDevice, error) {
	var device entities.NetwatchDevice
	err := r.db.Where("customer_id = ?", customerID).First(&device).Error
	return &device, err
}

// GetDeviceEvents returns events for a device within a time range
func (r *Repo) GetDeviceEvents(deviceID string, startTime, endTime time.Time) ([]entities.NetwatchEvent, error) {
	var events []entities.NetwatchEvent
	query := r.db.Where("device_id = ?", deviceID)
	
	if !startTime.IsZero() {
		query = query.Where("event_time >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("event_time <= ?", endTime)
	}
	
	err := query.Order("event_time DESC").Find(&events).Error
	return events, err
}

// GetLatestDeviceEvent returns the latest event for a device
func (r *Repo) GetLatestDeviceEvent(deviceID string) (*entities.NetwatchEvent, error) {
	var event entities.NetwatchEvent
	err := r.db.Where("device_id = ?", deviceID).
		Order("event_time DESC").
		First(&event).Error
	return &event, err
}

// GetCustomerInfo returns customer information with related data
func (r *Repo) GetCustomerInfo(customerID string) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.db.Preload("Product").Preload("Area").Preload("Company").
		Where("id = ?", customerID).First(&customer).Error
	return &customer, err
}

// GetDeviceStatusHistory returns status history for a device
func (r *Repo) GetDeviceStatusHistory(deviceID string, timeRange string) ([]entities.NetwatchEvent, error) {
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

	return r.GetDeviceEvents(deviceID, startTime, time.Time{})
}

// GetUptimeStats calculates uptime statistics for a device
func (r *Repo) GetUptimeStats(deviceID string, days int) (map[string]interface{}, error) {
	startTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	
	events, err := r.GetDeviceEvents(deviceID, startTime, time.Time{})
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_events":    len(events),
		"up_events":       0,
		"down_events":     0,
		"uptime_percent":  100.0,
		"total_downtime":  "0h 0m",
		"last_down_time":  nil,
		"last_up_time":    nil,
	}

	if len(events) == 0 {
		return stats, nil
	}

	// Count events
	for _, event := range events {
		if event.EventType == "up" {
			stats["up_events"] = stats["up_events"].(int) + 1
			if stats["last_up_time"] == nil {
				stats["last_up_time"] = event.EventTime
			}
		} else if event.EventType == "down" {
			stats["down_events"] = stats["down_events"].(int) + 1
			if stats["last_down_time"] == nil {
				stats["last_down_time"] = event.EventTime
			}
		}
	}

	// Calculate uptime percentage
	totalEvents := len(events)
	upEvents := stats["up_events"].(int)
	if totalEvents > 0 {
		stats["uptime_percent"] = float64(upEvents) / float64(totalEvents) * 100
	}

	return stats, nil
}
