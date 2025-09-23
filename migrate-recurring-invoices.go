package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

	// Execute first migration file
	fmt.Println("Creating recurring_invoices table...")
	migrationSQL1, err := ioutil.ReadFile("database-migration-recurring-invoices.sql")
	if err != nil {
		log.Fatalf("Error reading migration file 1: %v", err)
	}
	_, err = db.Exec(string(migrationSQL1))
	if err != nil {
		log.Fatalf("Error executing migration 1: %v", err)
	}
	fmt.Println("✓ recurring_invoices table created successfully")

	// Execute second migration file
	fmt.Println("Creating recurring_invoice_history table...")
	migrationSQL2, err := ioutil.ReadFile("database-migration-recurring-invoice-history.sql")
	if err != nil {
		log.Fatalf("Error reading migration file 2: %v", err)
	}
	_, err = db.Exec(string(migrationSQL2))
	if err != nil {
		log.Fatalf("Error executing migration 2: %v", err)
	}
	fmt.Println("✓ recurring_invoice_history table created successfully")

	// Execute index migration files
	fmt.Println("Creating indexes...")

	// Index 1
	migrationSQL3, err := ioutil.ReadFile("database-migration-recurring-invoice-indexes.sql")
	if err != nil {
		log.Fatalf("Error reading migration file 3: %v", err)
	}
	_, err = db.Exec(string(migrationSQL3))
	if err != nil {
		log.Fatalf("Error executing migration 3: %v", err)
	}

	// Index 2
	migrationSQL4, err := ioutil.ReadFile("database-migration-recurring-invoice-indexes-2.sql")
	if err != nil {
		log.Fatalf("Error reading migration file 4: %v", err)
	}
	_, err = db.Exec(string(migrationSQL4))
	if err != nil {
		log.Fatalf("Error executing migration 4: %v", err)
	}

	// Index 3
	migrationSQL5, err := ioutil.ReadFile("database-migration-recurring-invoice-indexes-3.sql")
	if err != nil {
		log.Fatalf("Error reading migration file 5: %v", err)
	}
	_, err = db.Exec(string(migrationSQL5))
	if err != nil {
		log.Fatalf("Error executing migration 5: %v", err)
	}

	fmt.Println("✓ indexes created successfully")

	fmt.Println("🎉 Recurring invoices migration completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
