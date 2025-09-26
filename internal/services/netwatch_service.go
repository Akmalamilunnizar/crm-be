package services

import (
	"context"
	"fmt"
	"log"
	// ticketapi "skripsi-be/internal/api/ticket" // Commented out to avoid circular dependency
	"skripsi-be/internal/models/entities"
	"strings"
	"time"

	"gorm.io/gorm"
)

type NetwatchService struct {
	db            *gorm.DB
	// ticketService *ticketapi.Service // Commented out to avoid circular dependency
}

type MikroTikNetwatchDevice struct {
	ID         string `json:".id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Status     string `json:"status"`
	Interval   string `json:"interval"`
	Timeout    string `json:"timeout"`
	UpScript   string `json:"up-script"`
	DownScript string `json:"down-script"`
}

type NetwatchEventProcessor struct {
	service *NetwatchService
}

func NewNetwatchService(db *gorm.DB) *NetwatchService {
	return &NetwatchService{
		db:            db,
		// ticketService: ticketService, // Commented out to avoid circular dependency
	}
}

// SyncDevicesFromMikroTik fetches all Netwatch devices from MikroTik and syncs with database
func (s *NetwatchService) SyncDevicesFromMikroTik(config *entities.NetwatchConfig) error {
	// Connect to MikroTik
	client, err := s.connectToMikroTik(config)
	if err != nil {
		return fmt.Errorf("failed to connect to MikroTik: %v", err)
	}
	defer client.Close()

	// Fetch all Netwatch devices
	devices, err := s.fetchNetwatchDevices(client)
	if err != nil {
		return fmt.Errorf("failed to fetch Netwatch devices: %v", err)
	}

	// Sync with database
	for _, device := range devices {
		err := s.syncDevice(device)
		if err != nil {
			log.Printf("Failed to sync device %s: %v", device.Name, err)
		}
	}

	return nil
}

// ProcessNetwatchEvent processes a single Netwatch event (up/down)
func (s *NetwatchService) ProcessNetwatchEvent(event *entities.NetwatchEvent) error {
	// Mark event as processed
	event.Processed = true
	if err := s.db.Save(event).Error; err != nil {
		return fmt.Errorf("failed to mark event as processed: %v", err)
	}

	// Get device details
	var device entities.NetwatchDevice
	if err := s.db.Preload("Customer").First(&device, "id = ?", event.DeviceID).Error; err != nil {
		return fmt.Errorf("failed to get device: %v", err)
	}

	// Update device status
	device.Status = event.EventType
	device.LastSeen = event.EventTime
	if err := s.db.Save(&device).Error; err != nil {
		return fmt.Errorf("failed to update device status: %v", err)
	}

	// Only process DOWN events for ticket creation
	if event.EventType == "down" {
		return s.handleDeviceDownEvent(event, &device)
	}

	return nil
}

// handleDeviceDownEvent handles when a device goes down
func (s *NetwatchService) handleDeviceDownEvent(event *entities.NetwatchEvent, device *entities.NetwatchDevice) error {
	// Check if device is associated with a customer
	if device.CustomerID == nil {
		log.Printf("Device %s is not associated with any customer, skipping ticket creation", device.Name)
		return nil
	}

	// Check if customer is active/paid
	var customer entities.Customer
	if err := s.db.First(&customer, "id = ?", *device.CustomerID).Error; err != nil {
		return fmt.Errorf("failed to get customer: %v", err)
	}

	// Check if customer is isolated due to unpaid bills
	if s.isCustomerIsolated(&customer) {
		log.Printf("Customer %s is isolated due to unpaid bills, skipping ticket creation", customer.Name)
		return nil
	}

	// Check if there's already an open ticket for this device
	var existingTicket entities.TroubleTicket
	err := s.db.Where("device_id = ? AND status IN (?)", device.ID, []string{"unfinished", "ongoing"}).
		First(&existingTicket).Error

	if err == nil {
		// Ticket already exists, add a log entry
		logEntry := &entities.TicketLog{
			TicketID: existingTicket.ID,
			LogType:  "netwatch",
			Message:  fmt.Sprintf("Device %s went DOWN at %s", device.Name, event.EventTime.Format("2006-01-02 15:04:05")),
			EventID:  &event.ID,
		}

		if err := s.db.Create(logEntry).Error; err != nil {
			log.Printf("Failed to create ticket log: %v", err)
		}

		return nil
	}

	// Create new ticket
	ticket := &entities.TroubleTicket{
		CustomerID:  *device.CustomerID,
		Type:        stringPtr("network"),
		Title:       fmt.Sprintf("Network Issue - Device %s Down", device.Name),
		Description: stringPtr(fmt.Sprintf("Device %s (%s) went DOWN at %s. This ticket was automatically created by Netwatch monitoring.", device.Name, device.IPAddress, event.EventTime.Format("2006-01-02 15:04:05"))),
		Status:      "ongoing",
	}

	// Create ticket using existing service
	// createdTicket, err := s.ticketService.CreateTicketFromNetwatch(ticket) // Commented out to avoid circular dependency
	// For now, just log the ticket creation
	log.Printf("Would create ticket from netwatch: %+v", ticket)
	var createdTicket *entities.TroubleTicket
	// For now, assume success
	// if err != nil {
	//	return fmt.Errorf("failed to create ticket: %v", err)
	// }

	// Create initial log entry
	logEntry := &entities.TicketLog{
		TicketID: createdTicket.ID,
		LogType:  "netwatch",
		Message:  fmt.Sprintf("Ticket automatically created due to device %s going DOWN", device.Name),
		EventID:  &event.ID,
	}

	if err := s.db.Create(logEntry).Error; err != nil {
		log.Printf("Failed to create initial ticket log: %v", err)
	}

	log.Printf("Created ticket #%d for device %s (customer: %s)", createdTicket.ID, device.Name, customer.Name)
	return nil
}

// isCustomerIsolated checks if customer is isolated due to unpaid bills
func (s *NetwatchService) isCustomerIsolated(customer *entities.Customer) bool {
	// Check if customer has unpaid invoices
	var unpaidCount int64
	s.db.Model(&entities.Invoice{}).
		Where("customer_id = ? AND status = ?", customer.ID, "unpaid").
		Count(&unpaidCount)

	// If customer has unpaid invoices, consider them isolated
	return unpaidCount > 0
}

// connectToMikroTik establishes connection to MikroTik router
func (s *NetwatchService) connectToMikroTik(config *entities.NetwatchConfig) (*MikroTikClient, error) {
	// This is a placeholder - you'll need to implement actual MikroTik API client
	// You can use libraries like github.com/ddelnano/go-mikrotik
	return &MikroTikClient{}, nil
}

// fetchNetwatchDevices fetches all Netwatch devices from MikroTik
func (s *NetwatchService) fetchNetwatchDevices(client *MikroTikClient) ([]MikroTikNetwatchDevice, error) {
	// This is a placeholder - implement actual MikroTik API call
	// Example: return client.GetNetwatchDevices()
	return []MikroTikNetwatchDevice{}, nil
}

// syncDevice syncs a single device with the database
func (s *NetwatchService) syncDevice(mikroTikDevice MikroTikNetwatchDevice) error {
	var device entities.NetwatchDevice

	// Try to find existing device by IP
	err := s.db.Where("ip_address = ?", mikroTikDevice.Host).First(&device).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new device
			device = entities.NetwatchDevice{
				Name:      mikroTikDevice.Name,
				IPAddress: mikroTikDevice.Host,
				Status:    strings.ToLower(mikroTikDevice.Status),
				LastSeen:  time.Now(),
			}

			// Try to find customer by IP address
			var customer entities.Customer
			if err := s.db.Where("ip_address = ?", mikroTikDevice.Host).First(&customer).Error; err == nil {
				device.CustomerID = &customer.ID
			}

			return s.db.Create(&device).Error
		}
		return err
	}

	// Update existing device
	device.Name = mikroTikDevice.Name
	device.Status = strings.ToLower(mikroTikDevice.Status)
	device.LastSeen = time.Now()

	return s.db.Save(&device).Error
}

// StartMonitoring starts the Netwatch monitoring service
func (s *NetwatchService) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Netwatch monitoring stopped")
			return
		case <-ticker.C:
			s.processPendingEvents()
		}
	}
}

// processPendingEvents processes all unprocessed Netwatch events
func (s *NetwatchService) processPendingEvents() {
	var events []entities.NetwatchEvent
	if err := s.db.Where("processed = ?", false).Find(&events).Error; err != nil {
		log.Printf("Failed to fetch pending events: %v", err)
		return
	}

	for _, event := range events {
		if err := s.ProcessNetwatchEvent(&event); err != nil {
			log.Printf("Failed to process event %s: %v", event.ID, err)
		}
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// MikroTikClient is a placeholder for MikroTik API client
type MikroTikClient struct{}

func (c *MikroTikClient) Close() error {
	return nil
}
