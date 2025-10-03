package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection from environment or defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "iqgncnzy_skripsi")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	fmt.Println("🚀 Starting Installation Technicians & MikroTik Provisioning Migration")
	fmt.Println("=" + strings.Repeat("=", 78))

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Error pinging database: %v", err)
	}

	fmt.Println("✅ Database connection successful")
	fmt.Printf("   Database: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	fmt.Println()

	// Read migration file
	migrationFile := "../migration_add_installation_technicians.sql"
	if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
		// Try alternative path
		migrationFile = "migration_add_installation_technicians.sql"
		if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
			log.Fatalf("❌ Migration file not found: %s", migrationFile)
		}
	}

	fmt.Printf("📄 Reading migration file: %s\n", migrationFile)
	sqlContent, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("❌ Failed to read migration file: %v", err)
	}

	fmt.Println("✅ Migration file loaded successfully")
	fmt.Println()

	// Execute migration steps
	fmt.Println("📋 Executing migration steps:")
	fmt.Println()

	steps := []struct {
		name        string
		description string
	}{
		{"Step 1", "Creating installation_report_technicians table"},
		{"Step 2", "Creating installation_provisioning_logs table"},
		{"Step 3", "Migrating existing technician data"},
		{"Step 4", "Adding provisioning columns to customer_installations"},
		{"Step 5", "Updating database views"},
		{"Step 6", "Creating triggers"},
		{"Step 7", "Finalizing migration"},
	}

	for _, step := range steps {
		fmt.Printf("⏳ %s: %s...\n", step.name, step.description)
	}
	fmt.Println()

	// Execute the entire migration
	fmt.Println("🔧 Executing migration...")
	_, err = db.Exec(string(sqlContent))
	if err != nil {
		log.Printf("❌ Migration failed: %v\n", err)
		fmt.Println()
		fmt.Println("💡 Troubleshooting tips:")
		fmt.Println("   1. Check if tables already exist")
		fmt.Println("   2. Verify database user has CREATE/ALTER privileges")
		fmt.Println("   3. Check for syntax errors in the migration file")
		fmt.Println("   4. Review error message above for specific issues")
		os.Exit(1)
	}

	fmt.Println("✅ Migration executed successfully!")
	fmt.Println()

	// Verify tables were created
	fmt.Println("🔍 Verifying migration...")

	tables := []string{
		"installation_report_technicians",
		"installation_provisioning_logs",
	}

	allTablesExist := true
	for _, table := range tables {
		var exists bool
		query := fmt.Sprintf("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s'", dbName, table)
		err := db.QueryRow(query).Scan(&exists)
		if err != nil || !exists {
			fmt.Printf("   ⚠️  Table '%s' not found\n", table)
			allTablesExist = false
		} else {
			fmt.Printf("   ✅ Table '%s' created\n", table)
		}
	}

	// Check new columns in customer_installations
	fmt.Println()
	fmt.Println("   Checking new columns in customer_installations:")

	columns := []string{
		"mac_address",
		"ip_address",
		"provisioning_status",
		"provisioning_completed_at",
		"code_name",
		"psb_date",
		"psb_time",
	}

	for _, column := range columns {
		var exists bool
		query := fmt.Sprintf("SELECT COUNT(*) > 0 FROM information_schema.columns WHERE table_schema = '%s' AND table_name = 'customer_installations' AND column_name = '%s'", dbName, column)
		err := db.QueryRow(query).Scan(&exists)
		if err != nil || !exists {
			fmt.Printf("   ⚠️  Column '%s' not added\n", column)
		} else {
			fmt.Printf("   ✅ Column '%s' added\n", column)
		}
	}

	// Count migrated records
	fmt.Println()
	fmt.Println("   Checking data migration:")

	var migratedCount int
	err = db.QueryRow("SELECT COUNT(*) FROM installation_report_technicians WHERE notes LIKE '%Migrated from legacy%'").Scan(&migratedCount)
	if err != nil {
		fmt.Printf("   ⚠️  Could not count migrated records: %v\n", err)
	} else {
		fmt.Printf("   ✅ Migrated %d existing technician assignments\n", migratedCount)
	}

	fmt.Println()
	if allTablesExist {
		fmt.Println("🎉 Migration completed successfully!")
	} else {
		fmt.Println("⚠️  Migration completed with warnings - please review above")
	}

	fmt.Println()
	fmt.Println("📋 Next steps:")
	fmt.Println("   1. ✅ Database schema updated")
	fmt.Println("   2. 🔄 Restart your Go application")
	fmt.Println("   3. 🧪 Test the new multiple technician feature")
	fmt.Println("   4. 🔧 Implement MikroTik connection in mikrotik_connection.go")
	fmt.Println("   5. 📝 Update API handlers as per MULTIPLE_TECHNICIAN_PROVISIONING_IMPLEMENTATION_GUIDE.md")
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 78))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
