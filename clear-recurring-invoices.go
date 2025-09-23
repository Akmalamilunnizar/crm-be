package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("🧹 Clearing Recurring Invoice Data...")

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

	// Clear recurring invoice history first (due to foreign key constraints)
	_, err = db.Exec("DELETE FROM recurring_invoice_history")
	if err != nil {
		log.Printf("⚠️  Error clearing recurring_invoice_history: %v", err)
	} else {
		fmt.Println("✅ Cleared recurring_invoice_history table")
	}

	// Clear recurring invoices
	_, err = db.Exec("DELETE FROM recurring_invoices")
	if err != nil {
		log.Printf("⚠️  Error clearing recurring_invoices: %v", err)
	} else {
		fmt.Println("✅ Cleared recurring_invoices table")
	}

	fmt.Println("🎉 Recurring invoice data cleared successfully!")
	fmt.Println("💡 You can now run 'go run seed-recurring-invoices.go' to reseed the data")
}
