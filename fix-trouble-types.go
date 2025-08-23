package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
)

func main() {
	db := database.GetDB()

	// First, let's see what's currently in the trouble_type table
	var existingTypes []entities.TroubleTypeRow
	if err := db.Table("trouble_type").Find(&existingTypes).Error; err != nil {
		log.Printf("Error querying existing trouble types: %v", err)
	} else {
		log.Printf("Existing trouble types: %+v", existingTypes)
	}

	// If no trouble types exist, create them
	if len(existingTypes) == 0 {
		log.Println("No trouble types found, creating default ones...")

		troubleTypes := []entities.TroubleTypeRow{
			{ID: uuid.New().String(), Name: stringPtr("Internet Down")},
			{ID: uuid.New().String(), Name: stringPtr("Slow Connection")},
			{ID: uuid.New().String(), Name: stringPtr("Modem Issue")},
			{ID: uuid.New().String(), Name: stringPtr("Cable Damage")},
			{ID: uuid.New().String(), Name: stringPtr("Power Outage")},
			{ID: uuid.New().String(), Name: stringPtr("Configuration Error")},
		}

		for _, tt := range troubleTypes {
			if err := db.Table("trouble_type").Create(&tt).Error; err != nil {
				log.Printf("Error creating trouble type %s: %v", *tt.Name, err)
			} else {
				log.Printf("Created trouble type: %s (ID: %s)", *tt.Name, tt.ID)
			}
		}
	}

	// Now let's check what types are used in tickets
	var ticketTypes []struct {
		ID   uint64  `json:"id"`
		Type *string `json:"type"`
	}
	if err := db.Table("trouble_tickets").Select("id, type").Find(&ticketTypes).Error; err != nil {
		log.Printf("Error querying ticket types: %v", err)
		return
	}
	log.Printf("Ticket types found: %+v", ticketTypes)

	// Get the first trouble type to use as default
	var defaultType entities.TroubleTypeRow
	if err := db.Table("trouble_type").First(&defaultType).Error; err != nil {
		log.Printf("Error getting default trouble type: %v", err)
		return
	}
	log.Printf("Default trouble type: %+v", defaultType)

	// Update tickets that have invalid or null type to use the default type
	result := db.Table("trouble_tickets").
		Where("type IS NULL OR type = '' OR type NOT IN (SELECT id FROM trouble_type)").
		Update("type", defaultType.ID)

	if result.Error != nil {
		log.Printf("Error updating tickets: %v", result.Error)
	} else {
		log.Printf("Updated %d tickets to use default trouble type", result.RowsAffected)
	}

	// Test the join again
	var result2 []struct {
		ID       uint64  `json:"id"`
		Type     *string `json:"type"`
		TypeName *string `json:"type_name"`
	}
	if err := db.Table("trouble_tickets t").
		Select("t.id, t.type, tt.name as type_name").
		Joins("LEFT JOIN trouble_type tt ON tt.id = t.type").
		Find(&result2).Error; err != nil {
		log.Printf("Error testing join: %v", err)
		return
	}
	log.Printf("Join result after fix: %+v", result2)
}

func stringPtr(s string) *string {
	return &s
}
