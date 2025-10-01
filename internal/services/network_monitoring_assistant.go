package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"skripsi-be/internal/models/entities"
	"gorm.io/gorm"
)

// NetworkMonitoringAssistant handles real-time device monitoring using MikroTik Netwatch
type NetworkMonitoringAssistant struct {
	db            *gorm.DB
	mikroTikService *MikroTikService
}

// DeviceStatus represents the current status of a monitored device
type DeviceStatus struct {
	IP           string    `json:"ip"`
	Status       string    `json:"status"` // "up" or "down"
	LastChecked  time.Time `json:"last_checked"`
	Recommendation string  `json:"recommendation,omitempty"`
}

// NewNetworkMonitoringAssistant creates a new instance
func NewNetworkMonitoringAssistant(db *gorm.DB, mikroTikService *MikroTikService) *NetworkMonitoringAssistant {
	return &NetworkMonitoringAssistant{
		db:            db,
		mikroTikService: mikroTikService,
	}
}

// MonitorCustomerDevice monitors a specific customer device by IP address
func (nma *NetworkMonitoringAssistant) MonitorCustomerDevice(customerID, ipAddress string) (*DeviceStatus, error) {
	// Check if device exists in netwatch
	var device entities.NetwatchDevice
	err := nma.db.Where("ip_address = ? AND customer_id = ?", ipAddress, customerID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Device not in netwatch, create it
			device = entities.NetwatchDevice{
				Name:      fmt.Sprintf("Customer-%s-Device", customerID),
				IPAddress: ipAddress,
				CustomerID: &customerID,
				Status:    "unknown",
				LastSeen:  time.Now(),
			}
			if err := nma.db.Create(&device).Error; err != nil {
				return nil, fmt.Errorf("failed to create netwatch device: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to query device: %v", err)
		}
	}

	// Check device status using MikroTik
	status, err := nma.checkDeviceStatus(ipAddress)
	if err != nil {
		log.Printf("Failed to check device status for %s: %v", ipAddress, err)
		status = "down"
	}

	// Update device status in database
	device.Status = status
	device.LastSeen = time.Now()
	nma.db.Save(&device)

	// Create netwatch event
	event := entities.NetwatchEvent{
		DeviceID:  device.ID,
		EventType: status,
		EventTime: time.Now(),
		RawData:   fmt.Sprintf("Manual check for IP %s", ipAddress),
	}
	nma.db.Create(&event)

	// Generate recommendation if device is down
	recommendation := ""
	if status == "down" {
		recommendation = nma.generateTroubleshootingRecommendation(ipAddress, &device)
	}

	return &DeviceStatus{
		IP:           ipAddress,
		Status:       status,
		LastChecked:  time.Now(),
		Recommendation: recommendation,
	}, nil
}

// checkDeviceStatus checks if a device is reachable using MikroTik ping
func (nma *NetworkMonitoringAssistant) checkDeviceStatus(ipAddress string) (string, error) {
	if nma.mikroTikService == nil || !nma.mikroTikService.IsConnected() {
		return "down", fmt.Errorf("MikroTik service not connected")
	}

	// Use MikroTik ping command to check device status
	pingCommand := fmt.Sprintf("/ping count=3 address=%s", ipAddress)
	output, err := nma.mikroTikService.ExecuteCommand(pingCommand)
	if err != nil {
		return "down", err
	}

	// Parse ping output to determine if device is up
	// MikroTik ping output typically shows packet loss percentage
	if strings.Contains(output, "0% packet loss") || strings.Contains(output, "received=3") {
		return "up", nil
	}

	return "down", nil
}

// generateTroubleshootingRecommendation provides troubleshooting tips for down devices
func (nma *NetworkMonitoringAssistant) generateTroubleshootingRecommendation(ipAddress string, device *entities.NetwatchDevice) string {
	recommendations := []string{
		"Periksa koneksi kabel LAN/WiFi",
		"Restart router/modem pelanggan",
		"Periksa status listrik di lokasi pelanggan",
		"Verifikasi konfigurasi IP address",
		"Periksa firewall atau security settings",
	}

	// Add specific recommendations based on device history
	if device != nil {
		// Check recent events for this device
		var recentEvents []entities.NetwatchEvent
		nma.db.Where("device_id = ? AND event_time > ?", device.ID, time.Now().Add(-24*time.Hour)).
			Order("event_time DESC").Limit(5).Find(&recentEvents)

		downCount := 0
		for _, event := range recentEvents {
			if event.EventType == "down" {
				downCount++
			}
		}

		if downCount >= 3 {
			recommendations = append(recommendations, "Device sering down - periksa hardware atau koneksi fisik")
		}
	}

	return strings.Join(recommendations, ", ")
}

// GetCustomerConnectionStatus returns the connection status for a customer
func (nma *NetworkMonitoringAssistant) GetCustomerConnectionStatus(customerID string) ([]DeviceStatus, error) {
	var devices []entities.NetwatchDevice
	err := nma.db.Where("customer_id = ?", customerID).Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get customer devices: %v", err)
	}

	var statuses []DeviceStatus
	for _, device := range devices {
		status := DeviceStatus{
			IP:          device.IPAddress,
			Status:      device.Status,
			LastChecked: device.LastSeen,
		}

		if device.Status == "down" {
			status.Recommendation = nma.generateTroubleshootingRecommendation(device.IPAddress, &device)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// GenerateNetwatchScript creates a MikroTik Netwatch script for monitoring an IP
func (nma *NetworkMonitoringAssistant) GenerateNetwatchScript(ipAddress, customerName string) string {
	// Sanitize customer name for script
	safeName := strings.ReplaceAll(customerName, " ", "_")
	safeName = strings.ReplaceAll(safeName, "-", "_")
	
	script := fmt.Sprintf(`# Netwatch Script for Customer: %s
# IP Address: %s
# Generated: %s

# Add netwatch entry
/ip firewall connection tracking set enabled=yes
/ip firewall connection tracking set tcp-established-timeout=1d
/ip firewall connection tracking set tcp-fin-wait-timeout=10s
/ip firewall connection tracking set tcp-close-wait-timeout=10s
/ip firewall connection tracking set tcp-syn-sent-timeout=5s
/ip firewall connection tracking set tcp-syn-received-timeout=5s
/ip firewall connection tracking set tcp-syn-ack-timeout=5s
/ip firewall connection tracking set tcp-last-ack-timeout=10s
/ip firewall connection tracking set tcp-time-wait-timeout=10s
/ip firewall connection tracking set tcp-close-timeout=10s
/ip firewall connection tracking set tcp-max-retrans-timeout=5m
/ip firewall connection tracking set tcp-unacked-timeout=5m
/ip firewall connection tracking set udp-timeout=10s
/ip firewall connection tracking set udp-stream-timeout=3m
/ip firewall connection tracking set icmp-timeout=10s
/ip firewall connection tracking set generic-timeout=10m

# Netwatch entry for customer device
/tool netwatch add host=%s interval=30s timeout=5s up-script=":log info \"Customer %s device UP: %s\"" down-script=":log warning \"Customer %s device DOWN: %s - Check connection\"" comment="Customer_%s_Device"

# Optional: Add email notification (uncomment if email is configured)
# /tool netwatch add host=%s interval=30s timeout=5s up-script=":log info \"Customer %s device UP: %s\"" down-script=":log warning \"Customer %s device DOWN: %s\"; /tool e-mail send to=\"admin@company.com\" subject=\"Device Down Alert\" body=\"Customer %s device (%s) is down. Please check connection.\"" comment="Customer_%s_Device_Email"

# Display current netwatch entries
:put "Netwatch entry added for customer: %s"
:put "IP Address: %s"
:put "Check status with: /tool netwatch print where host=%s"`,
		customerName,
		ipAddress,
		time.Now().Format("2006-01-02 15:04:05"),
		ipAddress,
		customerName,
		ipAddress,
		customerName,
		ipAddress,
		safeName,
		ipAddress,
		customerName,
		ipAddress,
		customerName,
		ipAddress,
		customerName,
		ipAddress,
		safeName,
		customerName,
		ipAddress,
		ipAddress,
	)

	return script
}

// FormatStatusResponse formats the status response for frontend display
func (nma *NetworkMonitoringAssistant) FormatStatusResponse(status *DeviceStatus) string {
	if status.Status == "up" {
		return fmt.Sprintf("[UP] %s – device online", status.IP)
	} else {
		return fmt.Sprintf("[DOWN] %s – device offline, rekomendasi: %s", status.IP, status.Recommendation)
	}
}

// StartRealTimeMonitoring starts real-time monitoring for all customer devices
func (nma *NetworkMonitoringAssistant) StartRealTimeMonitoring() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nma.checkAllDevices()
		}
	}
}

// checkAllDevices checks status of all monitored devices
func (nma *NetworkMonitoringAssistant) checkAllDevices() {
	var devices []entities.NetwatchDevice
	err := nma.db.Where("customer_id IS NOT NULL").Find(&devices).Error
	if err != nil {
		log.Printf("Failed to get devices for monitoring: %v", err)
		return
	}

	for _, device := range devices {
		status, err := nma.checkDeviceStatus(device.IPAddress)
		if err != nil {
			log.Printf("Failed to check device %s: %v", device.IPAddress, err)
			continue
		}

		// Update device status if changed
		if device.Status != status {
			device.Status = status
			device.LastSeen = time.Now()
			nma.db.Save(&device)

			// Create event
			event := entities.NetwatchEvent{
				DeviceID:  device.ID,
				EventType: status,
				EventTime: time.Now(),
				RawData:   fmt.Sprintf("Auto check - status changed to %s", status),
			}
			nma.db.Create(&event)

			log.Printf("Device %s status changed to %s", device.IPAddress, status)
		}
	}
}
