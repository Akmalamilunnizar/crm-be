package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = ""
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "iqgncnzy_skripsi"
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// Add due_date column
	fmt.Println("Adding due_date column to invoices table...")
	_, err = db.Exec("ALTER TABLE `invoices` ADD COLUMN `due_date` DATE NULL AFTER `status`")
	if err != nil {
		log.Printf("Warning: Column might already exist: %v", err)
	} else {
		fmt.Println("✅ Added due_date column")
	}

	// Add index
	fmt.Println("Adding index for due_date...")
	_, err = db.Exec("CREATE INDEX `idx_invoices_due_date` ON `invoices` (`due_date`)")
	if err != nil {
		log.Printf("Warning: Index might already exist: %v", err)
	} else {
		fmt.Println("✅ Added index for due_date")
	}

	fmt.Println("Migration completed successfully!")
}
