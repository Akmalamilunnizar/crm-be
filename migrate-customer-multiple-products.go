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

	// Read and execute migration file
	migrationFile := "database-migration-customer-multiple-products.sql"
	if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
		log.Fatalf("Migration file not found: %s", migrationFile)
	}

	sqlContent, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	// Split SQL content by semicolon and execute each statement
	statements := splitSQLStatements(string(sqlContent))

	for i, statement := range statements {
		statement = trimSQLStatement(statement)
		if statement == "" {
			continue
		}

		fmt.Printf("Executing statement %d...\n", i+1)
		fmt.Printf("SQL: %s\n", statement)

		_, err := db.Exec(statement)
		if err != nil {
			log.Printf("⚠️  Warning: Statement %d failed: %v", i+1, err)
			// Continue with other statements
		} else {
			fmt.Printf("✅ Statement %d executed successfully\n", i+1)
		}
		fmt.Println("---")
	}

	fmt.Println("🎉 Migration completed!")
	fmt.Println("\n📋 Next steps:")
	fmt.Println("1. Verify the migration worked correctly")
	fmt.Println("2. Test your application with the new schema")
	fmt.Println("3. Uncomment the final DROP statements in the migration file if everything works")
	fmt.Println("4. Run the migration again to remove the old product_id column from customer table")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inComment := false
	inString := false
	var stringChar rune

	for i, char := range sql {
		switch char {
		case '\'':
			if !inComment {
				if !inString {
					inString = true
					stringChar = '\''
				} else if stringChar == '\'' {
					inString = false
				}
			}
		case '"':
			if !inComment {
				if !inString {
					inString = true
					stringChar = '"'
				} else if stringChar == '"' {
					inString = false
				}
			}
		case '-':
			if !inString && !inComment && i+1 < len(sql) && sql[i+1] == '-' {
				inComment = true
			}
		case '\n':
			if inComment {
				inComment = false
			}
		case ';':
			if !inString && !inComment {
				statements = append(statements, current.String())
				current.Reset()
				continue
			}
		}

		if !inComment {
			current.WriteRune(char)
		}
	}

	// Add the last statement if it doesn't end with semicolon
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}

func trimSQLStatement(statement string) string {
	// Remove leading/trailing whitespace
	statement = strings.TrimSpace(statement)

	// Skip empty statements and comments
	if statement == "" || strings.HasPrefix(statement, "--") {
		return ""
	}

	return statement
}
