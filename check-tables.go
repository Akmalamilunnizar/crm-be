package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "iqgncnzy_skripsi:XhYJOWlwNgsk@tcp(103.63.24.139:3306)/iqgncnzy_skripsi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database successfully")
	log.Println("Checking for tables...")

	// Check if network_devices table exists
	var tables []string
	err = db.Raw("SHOW TABLES LIKE 'network%'").Scan(&tables).Error
	if err != nil {
		log.Fatal("Failed to query tables:", err)
	}

	fmt.Println("\nTables starting with 'network':")
	for _, table := range tables {
		fmt.Println("  -", table)
	}

	// Check all tables
	var allTables []string
	err = db.Raw("SHOW TABLES").Scan(&allTables).Error
	if err != nil {
		log.Fatal("Failed to query all tables:", err)
	}

	fmt.Println("\nAll tables in database:")
	for _, table := range allTables {
		fmt.Println("  -", table)
	}

	// Check if the specific tables we need exist
	requiredTables := []string{
		"customer_installations",
		"customer",
		"users",
		"network_devices",
		"assets",
		"products",
		"customer_services",
	}

	fmt.Println("\nChecking required tables:")
	for _, tableName := range requiredTables {
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", tableName).Scan(&count).Error
		if err != nil {
			fmt.Printf("  ❌ %s - Error: %v\n", tableName, err)
		} else if count > 0 {
			fmt.Printf("  ✅ %s - EXISTS\n", tableName)
		} else {
			fmt.Printf("  ❌ %s - NOT FOUND\n", tableName)
		}
	}
}
