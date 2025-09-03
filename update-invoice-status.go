package main

import (
	"fmt"
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
)

func main() {
	fmt.Println("Updating invoice status to pending...")

	// Get database connection
	db := database.GetDB()

	// Find the first unpaid invoice and update it to pending
	var invoice entities.Invoice
	result := db.Where("status = ?", entities.InvoiceStatusUnpaid).First(&invoice)
	if result.Error != nil {
		log.Fatalf("Error finding unpaid invoice: %v", result.Error)
	}

	// Update status to pending
	err := db.Model(&invoice).Update("status", entities.InvoiceStatusPending).Error
	if err != nil {
		log.Fatalf("Error updating invoice status: %v", err)
	}

	fmt.Printf("Updated invoice %s to pending status\n", invoice.ID)
	fmt.Printf("Invoice amount: %d IDR\n", invoice.Amount)
}