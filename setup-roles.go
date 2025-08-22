package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database connection
	db := database.GetDB()

	// Define the roles we need
	requiredRoles := []string{
		"ADMIN",
		"CUSTOMER_SERVICE",
		"NOC",
		"TECHNICIAN",
		"FINANCE",
	}

	// Check and create each role
	for _, roleName := range requiredRoles {
		var existingRole entities.Role
		if err := db.Where("name = ?", roleName).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			newRole := entities.Role{
				ID:   uuid.New().String(),
				Name: roleName,
			}
			if err := db.Create(&newRole).Error; err != nil {
				log.Printf("Error creating role %s: %v", roleName, err)
			} else {
				log.Printf("Created role: %s (ID: %s)", roleName, newRole.ID)
			}
		} else {
			log.Printf("Role %s already exists (ID: %s)", roleName, existingRole.ID)
		}
	}

	// List all roles in database
	var allRoles []entities.Role
	if err := db.Find(&allRoles).Error; err != nil {
		log.Printf("Error fetching roles: %v", err)
	} else {
		log.Println("All roles in database:")
		for _, role := range allRoles {
			log.Printf("  - %s (ID: %s)", role.Name, role.ID)
		}
	}

	log.Println("Role creation completed!")
}
