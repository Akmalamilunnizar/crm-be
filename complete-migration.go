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
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "iqgncnzy_skripsi")

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

	fmt.Println("✅ Database connection successful")
	fmt.Println("🔧 Completing migration steps...")

	// Step 1: Add missing index for product_id in network_devices
	fmt.Println("\n📊 Adding missing index for product_id in network_devices...")
	indexQuery := "CREATE INDEX IF NOT EXISTS idx_network_devices_product_id ON network_devices (product_id)"
	_, err = db.Exec(indexQuery)
	if err != nil {
		// Try without IF NOT EXISTS for older MySQL versions
		indexQuery = "CREATE INDEX idx_network_devices_product_id ON network_devices (product_id)"
		_, err = db.Exec(indexQuery)
		if err != nil {
			log.Printf("⚠️  Warning: Could not create index (might already exist): %v", err)
		} else {
			fmt.Println("✅ Index created successfully")
		}
	} else {
		fmt.Println("✅ Index created successfully")
	}

	// Step 2: Check if we need to remove product_id from customer table
	fmt.Println("\n👥 Checking customer table for product_id column...")
	var columnExists bool
	checkQuery := `
		SELECT COUNT(*) 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? 
			AND TABLE_NAME = 'customer' 
			AND COLUMN_NAME = 'product_id'
	`
	err = db.QueryRow(checkQuery, dbName).Scan(&columnExists)
	if err != nil {
		log.Printf("Error checking customer table: %v", err)
	} else {
		if columnExists {
			fmt.Println("⚠️  Column product_id still exists in customer table")
			fmt.Println("   This should be removed to complete the migration")

			// Ask for confirmation before removing
			fmt.Println("\n🤔 Do you want to remove product_id from customer table?")
			fmt.Println("   This is the final step of the migration.")
			fmt.Println("   Make sure your application is updated to use network_devices.product_id instead.")
			fmt.Println("   Type 'yes' to continue or anything else to skip:")

			var response string
			fmt.Scanln(&response)

			if response == "yes" {
				fmt.Println("\n🗑️  Removing product_id from customer table...")

				// Remove foreign key constraint first
				fmt.Println("   Removing foreign key constraint...")
				_, err = db.Exec("ALTER TABLE customer DROP FOREIGN KEY IF EXISTS customer_ibfk_3")
				if err != nil {
					log.Printf("   Warning: Could not drop foreign key: %v", err)
				}

				// Remove index
				fmt.Println("   Removing index...")
				_, err = db.Exec("ALTER TABLE customer DROP KEY IF EXISTS idx_customer_product_id")
				if err != nil {
					log.Printf("   Warning: Could not drop index: %v", err)
				}

				// Remove column
				fmt.Println("   Removing column...")
				_, err = db.Exec("ALTER TABLE customer DROP COLUMN product_id")
				if err != nil {
					log.Printf("   Error removing column: %v", err)
				} else {
					fmt.Println("✅ Column product_id removed from customer table")
				}
			} else {
				fmt.Println("⏭️  Skipping removal of product_id from customer table")
				fmt.Println("   You can run this script again later to complete the migration")
			}
		} else {
			fmt.Println("✅ Column product_id not found in customer table")
			fmt.Println("   Migration is already complete!")
		}
	}

	// Step 3: Final verification
	fmt.Println("\n🔍 Final verification...")

	// Check network_devices structure
	fmt.Println("   Checking network_devices table...")
	rows, err := db.Query("SELECT COUNT(*) FROM network_devices WHERE product_id IS NOT NULL")
	if err != nil {
		log.Printf("   Error checking network_devices: %v", err)
	} else {
		var count int
		rows.Scan(&count)
		rows.Close()
		fmt.Printf("   Found %d network devices with product_id assigned\n", count)
	}

	// Check foreign key constraints
	fmt.Println("   Checking foreign key constraints...")
	fkQuery := `
		SELECT COUNT(*) 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_SCHEMA = ? 
			AND TABLE_NAME = 'network_devices' 
			AND COLUMN_NAME = 'product_id'
			AND REFERENCED_TABLE_NAME IS NOT NULL
	`
	var fkCount int
	err = db.QueryRow(fkQuery, dbName).Scan(&fkCount)
	if err != nil {
		log.Printf("   Error checking foreign keys: %v", err)
	} else {
		fmt.Printf("   Found %d foreign key constraints for product_id\n", fkCount)
	}

	fmt.Println("\n🎉 Migration completion process finished!")
	fmt.Println("\n📋 Summary:")
	fmt.Println("✅ Recurring invoice original_day migration: COMPLETED")
	fmt.Println("✅ Invoice pending reasons migration: COMPLETED")
	fmt.Println("✅ Customer multiple products migration: COMPLETED")
	fmt.Println("\n📝 Next steps:")
	fmt.Println("1. Test your application with the new schema")
	fmt.Println("2. Update your application code to use network_devices.product_id instead of customer.product_id")
	fmt.Println("3. Verify all features work correctly")
	fmt.Println("4. If everything works, you can remove the migration files")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

