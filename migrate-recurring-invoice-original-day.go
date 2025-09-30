package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Get database connection string from environment or use default
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "iqgncnzy_skripsi")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Execute migration statements one by one
	fmt.Println("Adding original_day column to recurring_invoices table...")

	// Statement 1: Add column
	_, err = db.Exec("ALTER TABLE `recurring_invoices` ADD COLUMN `original_day` int NOT NULL DEFAULT 1 COMMENT 'Preserve the original template day (1-31) for recurring invoice calculations'")
	if err != nil {
		log.Fatalf("Error adding original_day column: %v", err)
	}
	fmt.Println("✓ original_day column added successfully")

	// Statement 2: Update existing records
	fmt.Println("Updating existing records...")
	_, err = db.Exec("UPDATE `recurring_invoices` SET `original_day` = DAY(`invoice_date`)")
	if err != nil {
		log.Fatalf("Error updating existing records: %v", err)
	}
	fmt.Println("✓ existing records updated successfully")

	// Statement 3: Add index
	fmt.Println("Adding index...")
	_, err = db.Exec("CREATE INDEX `idx_recurring_invoices_original_day` ON `recurring_invoices` (`original_day`)")
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}
	fmt.Println("✓ index created successfully")

	fmt.Println("🎉 Recurring invoice original_day migration completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
