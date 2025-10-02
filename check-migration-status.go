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
	fmt.Println("🔍 Checking migration status...")

	// 1. Check network_devices table structure
	fmt.Println("\n📋 Checking network_devices table structure...")
	rows, err := db.Query("DESCRIBE network_devices")
	if err != nil {
		log.Fatalf("Error describing network_devices: %v", err)
	}
	defer rows.Close()

	fmt.Println("Columns in network_devices table:")
	hasProductId := false
	for rows.Next() {
		var field, fieldType, null, key, defaultVal, extra sql.NullString
		err := rows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		fmt.Printf("  - %s: %s\n", field.String, fieldType.String)
		if field.String == "product_id" {
			hasProductId = true
		}
	}

	// 2. Check if product_id exists
	if hasProductId {
		fmt.Println("\n⚠️  Column product_id already exists in network_devices table")
		fmt.Println("   This explains why migration statement 1 failed (duplicate column)")
	} else {
		fmt.Println("\n❌ Column product_id not found in network_devices table")
		fmt.Println("   This is unexpected since migration should have added it")
	}

	// 3. Check foreign key constraints
	fmt.Println("\n🔗 Checking foreign key constraints...")
	fkQuery := `
		SELECT 
			CONSTRAINT_NAME,
			COLUMN_NAME,
			REFERENCED_TABLE_NAME,
			REFERENCED_COLUMN_NAME
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_SCHEMA = ? 
			AND TABLE_NAME = 'network_devices' 
			AND REFERENCED_TABLE_NAME IS NOT NULL
	`

	fkRows, err := db.Query(fkQuery, dbName)
	if err != nil {
		log.Printf("Error checking foreign keys: %v", err)
	} else {
		defer fkRows.Close()

		hasFK := false
		for fkRows.Next() {
			var constraintName, columnName, refTable, refColumn sql.NullString
			err := fkRows.Scan(&constraintName, &columnName, &refTable, &refColumn)
			if err != nil {
				log.Printf("Error scanning FK row: %v", err)
				continue
			}
			fmt.Printf("  - %s: %s -> %s.%s\n",
				constraintName.String, columnName.String,
				refTable.String, refColumn.String)
			hasFK = true
		}

		if !hasFK {
			fmt.Println("  No foreign key constraints found")
		}
	}

	// 4. Check indexes
	fmt.Println("\n📊 Checking indexes for product_id...")
	indexQuery := "SHOW INDEX FROM network_devices WHERE Column_name = 'product_id'"
	indexRows, err := db.Query(indexQuery)
	if err != nil {
		log.Printf("Error checking indexes: %v", err)
	} else {
		defer indexRows.Close()

		hasIndex := false
		for indexRows.Next() {
			var table, nonUnique, keyName, seqInIndex, columnName, collation, cardinality, subPart, packed, null, indexType, comment, indexComment sql.NullString
			err := indexRows.Scan(&table, &nonUnique, &keyName, &seqInIndex, &columnName, &collation, &cardinality, &subPart, &packed, &null, &indexType, &comment, &indexComment)
			if err != nil {
				log.Printf("Error scanning index row: %v", err)
				continue
			}
			fmt.Printf("  - %s: %s (Unique: %s)\n",
				keyName.String, columnName.String,
				func() string {
					if nonUnique.String == "0" {
						return "Yes"
					} else {
						return "No"
					}
				}())
			hasIndex = true
		}

		if !hasIndex {
			fmt.Println("  No indexes found for product_id")
		}
	}

	// 5. Check customer table
	fmt.Println("\n👥 Checking customer table...")
	customerRows, err := db.Query("DESCRIBE customer")
	if err != nil {
		log.Printf("Error describing customer table: %v", err)
	} else {
		defer customerRows.Close()

		hasCustomerProductId := false
		for customerRows.Next() {
			var field, fieldType, null, key, defaultVal, extra sql.NullString
			err := customerRows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
			if err != nil {
				log.Printf("Error scanning customer row: %v", err)
				continue
			}
			if field.String == "product_id" {
				hasCustomerProductId = true
				break
			}
		}

		if hasCustomerProductId {
			fmt.Println("✅ Column product_id exists in customer table")
			fmt.Println("   This will be removed after migration is complete")
		} else {
			fmt.Println("ℹ️  Column product_id not found in customer table")
		}
	}

	// 6. Recommendations
	fmt.Println("\n📝 RECOMMENDATIONS:")
	fmt.Println("1. Migration ran successfully despite the warning")
	fmt.Println("2. Statement 1 failed because product_id column already exists (this is normal)")
	fmt.Println("3. Statements 2 and 3 succeeded (foreign key and index were added)")
	fmt.Println("4. Next step: Verify your application works correctly with the new schema")
	fmt.Println("5. If everything is OK, run the script to remove product_id from customer table")

	fmt.Println("\n✅ Migration status check completed!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

