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

	// Create MikroTik configuration
	config := entities.NetwatchConfig{
		ID:       uuid.New().String(),
		Name:     "Main MikroTik Router",
		Host:     "192.168.1.1", // Ganti dengan IP MikroTik yang benar
		Port:     8728,
		Username: "admin", // Ganti dengan username yang benar
		Password: "admin", // Ganti dengan password yang benar
		UseSSL:   false,
		Active:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if config already exists
	var existingConfig entities.NetwatchConfig
	err := db.Where("active = ?", true).First(&existingConfig).Error
	if err == nil {
		fmt.Printf("Active config already exists: %s\n", existingConfig.ID)
		// Update existing config
		existingConfig.Host = config.Host
		existingConfig.Port = config.Port
		existingConfig.Username = config.Username
		existingConfig.Password = config.Password
		existingConfig.UpdatedAt = time.Now()
		err = db.Save(&existingConfig).Error
		if err != nil {
			log.Printf("Failed to update config: %v", err)
			return
		}
		fmt.Printf("Updated existing config\n")
	} else {
		// Create new config
		err = db.Create(&config).Error
		if err != nil {
			log.Printf("Failed to create config: %v", err)
			return
		}
		fmt.Printf("Created new config: %s\n", config.ID)
	}

	// Verify the config
	var configs []entities.NetwatchConfig
	err = db.Where("active = ?", true).Find(&configs).Error
	if err != nil {
		log.Printf("Failed to query configs: %v", err)
		return
	}

	fmt.Printf("\nActive MikroTik Configurations:\n")
	for _, c := range configs {
		fmt.Printf("ID: %s, Name: %s, Host: %s:%d, Username: %s\n", 
			c.ID, c.Name, c.Host, c.Port, c.Username)
	}

	fmt.Println("\nMikroTik configuration added successfully!")
	fmt.Println("IMPORTANT: Update the Host, Username, and Password with your actual MikroTik credentials!")
}
