package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
)

func main() {
	db := database.GetDB()

	// Check trouble_type table
	var troubleTypes []entities.TroubleTypeRow
	if err := db.Table("trouble_type").Find(&troubleTypes).Error; err != nil {
		log.Printf("Error querying trouble_type: %v", err)
		return
	}
	log.Printf("Trouble types in database: %+v", troubleTypes)

	// Check what types are used in tickets
	var ticketTypes []struct {
		ID   uint64  `json:"id"`
		Type *string `json:"type"`
	}
	if err := db.Table("trouble_tickets").Select("id, type").Find(&ticketTypes).Error; err != nil {
		log.Printf("Error querying trouble_tickets: %v", err)
		return
	}
	log.Printf("Ticket types: %+v", ticketTypes)

	// Test the join
	var result []struct {
		ID       uint64  `json:"id"`
		Type     *string `json:"type"`
		TypeName *string `json:"type_name"`
	}
	if err := db.Table("trouble_tickets t").
		Select("t.id, t.type, tt.name as type_name").
		Joins("LEFT JOIN trouble_type tt ON tt.id = t.type").
		Find(&result).Error; err != nil {
		log.Printf("Error testing join: %v", err)
		return
	}
	log.Printf("Join result: %+v", result)
}
