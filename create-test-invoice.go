package main

import (
	"fmt"
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
)

func main() {
	fmt.Println("Creating test invoice for partial payment...")

	// Get database connection
	db := database.GetDB()

	// Create a test invoice with pending status
	testInvoice := entities.Invoice{
		Amount:     500000, // 500,000 IDR
		CustomerID: "24225552-b0c3-43d7-9a67-20e04f36fa5f", // Using existing customer
		Link:       "",
		Status:     entities.InvoiceStatusPending,
	}

	// Create the invoice
	err := db.Create(&testInvoice).Error
	if err != nil {
		log.Fatalf("Error creating test invoice: %v", err)
	}

	// Update status to pending after creation to override BeforeCreate hook
	err = db.Model(&testInvoice).Update("status", entities.InvoiceStatusPending).Error
	if err != nil {
		log.Printf("Error updating invoice status: %v", err)
	}

	fmt.Printf("Created test invoice with ID: %s\n", testInvoice.ID)

	// Create some invoice items
	testItems := []entities.InvoiceItems{
		{
			Name:      "Internet Package 25 Mbps",
			Qty:       1,
			Price:     350000,
			Total:     350000,
			InvoiceID: testInvoice.ID,
		},
		{
			Name:      "Installation Fee",
			Qty:       1,
			Price:     150000,
			Total:     150000,
			InvoiceID: testInvoice.ID,
		},
	}

	for _, item := range testItems {
		err := db.Create(&item).Error
		if err != nil {
			log.Printf("Error creating invoice item: %v", err)
		} else {
			fmt.Printf("Created invoice item: %s\n", item.Name)
		}
	}

	fmt.Println("Test invoice creation completed!")
	fmt.Printf("Invoice ID: %s\n", testInvoice.ID)
	fmt.Printf("Total Amount: %d IDR\n", testInvoice.Amount)
	fmt.Printf("Status: %s\n", testInvoice.Status)
}