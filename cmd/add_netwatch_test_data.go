package main

import (
	"fmt"
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Initialize database
	db := database.GetDB()

	// Find customer with name 'sadvwabafbsdgd' or any customer
	var customer entities.Customer
	err := db.Where("name = ?", "sadvwabafbsdgd").First(&customer).Error
	if err != nil {
		// Try to find any customer
		err = db.First(&customer).Error
		if err != nil {
			log.Printf("No customers found: %v", err)
			return
		}
		log.Printf("Using customer: %s (ID: %s)", customer.Name, customer.ID)
	}

	fmt.Printf("Found customer: %s (ID: %s)\n", customer.Name, customer.ID)

	// Create netwatch device
	device := entities.NetwatchDevice{
		ID:         uuid.New().String(),
		Name:       "Customer-sadvwabafbsdgd-Device",
		IPAddress:  "10.10.21.10",
		CustomerID: &customer.ID,
		Status:     "up",
		LastSeen:   time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Check if device already exists
	var existingDevice entities.NetwatchDevice
	err = db.Where("ip_address = ?", "10.10.21.10").First(&existingDevice).Error
	if err == nil {
		fmt.Printf("Device already exists: %s\n", existingDevice.ID)
		// Update existing device
		existingDevice.Status = "up"
		existingDevice.LastSeen = time.Now()
		existingDevice.UpdatedAt = time.Now()
		err = db.Save(&existingDevice).Error
		if err != nil {
			log.Printf("Failed to update device: %v", err)
			return
		}
		fmt.Printf("Updated existing device\n")
	} else {
		// Create new device
		err = db.Create(&device).Error
		if err != nil {
			log.Printf("Failed to create device: %v", err)
			return
		}
		fmt.Printf("Created new device: %s\n", device.ID)
	}

	// Create netwatch event
	event := entities.NetwatchEvent{
		ID:        uuid.New().String(),
		DeviceID:  device.ID,
		EventType: "up",
		EventTime: time.Now(),
		RawData:   "Test event - device UP",
		Processed: true,
		CreatedAt: time.Now(),
	}

	err = db.Create(&event).Error
	if err != nil {
		log.Printf("Failed to create event: %v", err)
		return
	}

	fmt.Printf("Created event: %s\n", event.ID)

	// Verify the data
	var devices []entities.NetwatchDevice
	err = db.Preload("Customer").Where("ip_address = ?", "10.10.21.10").Find(&devices).Error
	if err != nil {
		log.Printf("Failed to query devices: %v", err)
		return
	}

	fmt.Printf("\nVerification:\n")
	for _, d := range devices {
		fmt.Printf("Device: %s, IP: %s, Status: %s, Customer: %s\n", 
			d.Name, d.IPAddress, d.Status, 
			func() string {
				if d.Customer != nil {
					return d.Customer.Name
				}
				return "No customer"
			}())
	}

	fmt.Println("\nTest data added successfully!")
}
