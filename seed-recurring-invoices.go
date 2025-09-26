package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RecurringInvoice struct {
	ID              string    `json:"id"`
	CustomerID      string    `json:"customer_id"`
	Amount          int       `json:"amount"`
	InvoiceDate     time.Time `json:"invoice_date"`
	DueDate         time.Time `json:"due_date"`
	NextInvoiceDate time.Time `json:"next_invoice_date"`
	Frequency       string    `json:"frequency"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	InvoiceItems    string    `json:"invoice_items"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       string    `json:"created_by"`
}

func main() {
	fmt.Println("🌱 Starting Recurring Invoice Seeder...")

	// Database connection
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/iqgncnzy_skripsi")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to database successfully!")

	// Check if recurring invoices already exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM recurring_invoices").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check existing recurring invoices:", err)
	}

	if count > 0 {
		fmt.Printf("⚠️  Found %d existing recurring invoices. Skipping seeding to avoid duplicates.\n", count)
		fmt.Println("💡 If you want to reseed, please clear the recurring_invoices table first.")
		return
	}

	// Get existing customer IDs
	customers := make(map[string]string)
	rows, err := db.Query("SELECT id, name FROM customer LIMIT 10")
	if err != nil {
		log.Fatal("Failed to query customers:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Failed to scan customer:", err)
		}
		customers[id] = name
	}

	if len(customers) == 0 {
		log.Fatal("No customers found in database. Please seed customers first.")
	}

	// Get existing user IDs
	users := make(map[string]string)
	rows2, err := db.Query("SELECT id, name FROM users LIMIT 5")
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id, name string
		if err := rows2.Scan(&id, &name); err != nil {
			log.Fatal("Failed to scan user:", err)
		}
		users[id] = name
	}

	if len(users) == 0 {
		log.Fatal("No users found in database. Please seed users first.")
	}

	// Get first customer and user IDs for seeding
	var firstCustomerID, firstUserID string
	for customerID := range customers {
		firstCustomerID = customerID
		break
	}
	for userID := range users {
		firstUserID = userID
		break
	}

	fmt.Printf("📋 Using customer: %s (%s)\n", firstCustomerID, customers[firstCustomerID])
	fmt.Printf("👤 Using user: %s (%s)\n", firstUserID, users[firstUserID])

	// Sample recurring invoices data
	recurringInvoices := []RecurringInvoice{
		{
			ID:              "rec_001",
			CustomerID:      firstCustomerID,
			Amount:          5000000,
			InvoiceDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			DueDate:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NextInvoiceDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Frequency:       "monthly",
			Status:          "active",
			Description:     "Monthly internet service subscription",
			InvoiceItems:    `[{"description":"Internet Service - 100 Mbps","quantity":1,"price":5000000,"total":5000000}]`,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CreatedBy:       firstUserID,
		},
		{
			ID:              "rec_002",
			CustomerID:      firstCustomerID,
			Amount:          2500000,
			InvoiceDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			DueDate:         time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			NextInvoiceDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Frequency:       "monthly",
			Status:          "active",
			Description:     "Monthly hosting and maintenance",
			InvoiceItems:    `[{"description":"Web Hosting - Premium Plan","quantity":1,"price":1500000,"total":1500000},{"description":"Maintenance Service","quantity":1,"price":1000000,"total":1000000}]`,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CreatedBy:       firstUserID,
		},
		{
			ID:              "rec_003",
			CustomerID:      firstCustomerID,
			Amount:          1000000,
			InvoiceDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			DueDate:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NextInvoiceDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Frequency:       "quarterly",
			Status:          "active",
			Description:     "Quarterly software license",
			InvoiceItems:    `[{"description":"Software License - Business Edition","quantity":1,"price":1000000,"total":1000000}]`,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CreatedBy:       firstUserID,
		},
		{
			ID:              "rec_004",
			CustomerID:      firstCustomerID,
			Amount:          7500000,
			InvoiceDate:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			DueDate:         time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			NextInvoiceDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Frequency:       "monthly",
			Status:          "active",
			Description:     "Monthly cloud services and support",
			InvoiceItems:    `[{"description":"Cloud Storage - 1TB","quantity":1,"price":2000000,"total":2000000},{"description":"Cloud Computing - Standard","quantity":1,"price":3000000,"total":3000000},{"description":"24/7 Support","quantity":1,"price":2500000,"total":2500000}]`,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CreatedBy:       firstUserID,
		},
		{
			ID:              "rec_005",
			CustomerID:      firstCustomerID,
			Amount:          3000000,
			InvoiceDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			DueDate:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			NextInvoiceDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			Frequency:       "monthly",
			Status:          "stopped",
			Description:     "Monthly digital marketing services",
			InvoiceItems:    `[{"description":"Social Media Management","quantity":1,"price":1500000,"total":1500000},{"description":"SEO Services","quantity":1,"price":1500000,"total":1500000}]`,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CreatedBy:       firstUserID,
		},
	}

	// Insert recurring invoices
	successCount := 0
	for _, invoice := range recurringInvoices {
		query := `
			INSERT INTO recurring_invoices (
				id, customer_id, amount, invoice_date, due_date, 
				next_invoice_date, frequency, status, description, 
				invoice_items, created_at, updated_at, created_by
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err := db.Exec(query,
			invoice.ID,
			invoice.CustomerID,
			invoice.Amount,
			invoice.InvoiceDate,
			invoice.DueDate,
			invoice.NextInvoiceDate,
			invoice.Frequency,
			invoice.Status,
			invoice.Description,
			invoice.InvoiceItems,
			invoice.CreatedAt,
			invoice.UpdatedAt,
			invoice.CreatedBy,
		)

		if err != nil {
			log.Printf("❌ Error inserting recurring invoice %s: %v", invoice.ID, err)
		} else {
			fmt.Printf("✅ Inserted recurring invoice: %s - %s (Rp %,.0f)\n",
				invoice.ID, invoice.Description, float64(invoice.Amount))
			successCount++
		}
	}

	fmt.Println("\n🎉 Recurring invoice seeding completed!")
	fmt.Printf("📊 Summary: %d/%d recurring invoices created successfully\n", successCount, len(recurringInvoices))

	if successCount > 0 {
		fmt.Println("💡 You can now test the recurring invoice feature in your application!")
		fmt.Println("🔗 Navigate to /dashboard/recurring-invoice to see the seeded data")
	}
}
