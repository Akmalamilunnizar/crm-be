package main

import (
	"fmt"
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
)

func main() {
	fmt.Println("Fixing Invoice Amounts...")

	// Get database connection
	db := database.GetDB()

	// Get all invoices with their items
	var invoices []entities.Invoice
	err := db.Preload("InvoiceItems").Find(&invoices).Error
	if err != nil {
		log.Fatalf("Error fetching invoices: %v", err)
	}

	fmt.Printf("Found %d invoices to process\n", len(invoices))

	// Process each invoice
	for i, invoice := range invoices {
		fmt.Printf("Processing invoice %d: %s\n", i+1, invoice.ID)

		if len(invoice.InvoiceItems) == 0 {
			fmt.Printf("  No invoice items found, skipping\n")
			continue
		}

		// Calculate total from items
		var totalAmount int64 = 0
		for _, item := range invoice.InvoiceItems {
			itemTotal := item.Price * item.Qty
			totalAmount += itemTotal
			fmt.Printf("  Item: %s, Price: %d, Qty: %d, Total: %d\n",
				item.Name, item.Price, item.Qty, itemTotal)
		}

		fmt.Printf("  Calculated total: %d, Current amount: %d\n", totalAmount, invoice.Amount)

		// Update if different
		if invoice.Amount != totalAmount {
			fmt.Printf("  Updating invoice amount from %d to %d\n", invoice.Amount, totalAmount)

			err := db.Model(&invoice).Update("amount", totalAmount).Error
			if err != nil {
				log.Printf("  Error updating invoice: %v", err)
			} else {
				fmt.Printf("  Successfully updated invoice amount\n")
			}
		} else {
			fmt.Printf("  Amount already correct, no update needed\n")
		}

		fmt.Println()
	}

	fmt.Println("Invoice amount fixing completed!")
}
