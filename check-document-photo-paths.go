package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CustomerInstallation struct {
	ID            string  `gorm:"column:id"`
	DocumentType  *string `gorm:"column:document_type"`
	DocumentPhoto *string `gorm:"column:document_photo"`
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("MYSQL")
	if dsn == "" {
		log.Fatal("MYSQL environment variable not set")
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🔍 Checking document photo paths in database...")
	fmt.Println("================================================")

	var installations []CustomerInstallation
	err = db.Where("document_photo IS NOT NULL AND document_photo != ''").Order("createdAt DESC").Limit(10).Find(&installations).Error
	if err != nil {
		log.Fatal("Failed to query database:", err)
	}

	if len(installations) == 0 {
		fmt.Println("❌ No document photos found in database")
		return
	}

	fmt.Printf("✅ Found %d document photo records:\n\n", len(installations))

	for i, installation := range installations {
		fmt.Printf("%d. ID: %s\n", i+1, installation.ID)
		if installation.DocumentType != nil {
			fmt.Printf("   Document Type: %s\n", *installation.DocumentType)
		}
		if installation.DocumentPhoto != nil {
			path := *installation.DocumentPhoto
			fmt.Printf("   Document Photo: %s\n", path)

			// Analyze path format
			if strings.Contains(path, "uploads/installations/documents/uploads/installations/documents/") {
				fmt.Printf("   ❌ Status: TRIPLE DUPLICATION\n")
			} else if strings.Contains(path, "uploads/installations/documents/uploads/installations/") {
				fmt.Printf("   ❌ Status: DOUBLE DUPLICATION\n")
			} else if strings.Contains(path, "uploads/installations/documents/uploads/") {
				fmt.Printf("   ❌ Status: SINGLE DUPLICATION\n")
			} else if strings.HasPrefix(path, "uploads/documents/") {
				fmt.Printf("   ❌ Status: WRONG STRUCTURE\n")
			} else if strings.HasPrefix(path, "uploads/installations/documents/") {
				fmt.Printf("   ✅ Status: CORRECT\n")
			} else {
				fmt.Printf("   ❓ Status: UNKNOWN FORMAT\n")
			}

			// Show what it should be normalized to
			normalized := normalizeDocumentPhotoPath(path)
			if normalized != path {
				fmt.Printf("   🔄 Should be: %s\n", normalized)
			}
		}
		fmt.Println()
	}

	fmt.Println("📋 Summary:")
	fmt.Println("===============")

	// Count different path types
	var correct, duplicated, wrongStructure int

	for _, installation := range installations {
		if installation.DocumentPhoto != nil {
			path := *installation.DocumentPhoto
			// Convert backslashes to forward slashes for consistent checking
			normalizedPath := strings.ReplaceAll(path, "\\", "/")

			if strings.HasPrefix(normalizedPath, "uploads/installations/documents/") && !strings.Contains(normalizedPath, "uploads/installations/documents/uploads/") {
				correct++
			} else if strings.Contains(normalizedPath, "uploads/installations/documents/uploads/") {
				duplicated++
			} else if strings.HasPrefix(normalizedPath, "uploads/documents/") {
				wrongStructure++
			}
		}
	}

	fmt.Printf("✅ Correct paths: %d\n", correct)
	fmt.Printf("❌ Duplicated paths: %d\n", duplicated)
	fmt.Printf("❌ Wrong structure: %d\n", wrongStructure)
}

// normalizeDocumentPhotoPath - Normalize document photo path
func normalizeDocumentPhotoPath(path string) string {
	// First, convert Windows backslashes to forward slashes for web URLs
	path = strings.ReplaceAll(path, "\\", "/")

	// Handle triple duplication
	for strings.Contains(path, "uploads/installations/documents/uploads/installations/documents/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/installations/documents/", "uploads/installations/documents/", 1)
	}

	// Handle double duplication
	for strings.Contains(path, "uploads/installations/documents/uploads/installations/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/installations/", "uploads/installations/documents/", 1)
	}

	// Handle single duplication
	if strings.HasPrefix(path, "uploads/installations/documents/uploads/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/", "uploads/installations/documents/", 1)
	}

	// Handle wrong structure
	if strings.HasPrefix(path, "uploads/documents/") {
		path = strings.Replace(path, "uploads/documents/", "uploads/installations/documents/", 1)
	}

	// Handle just filenames
	if !strings.Contains(path, "/") && strings.HasSuffix(path, ".jpg") {
		path = "uploads/installations/documents/" + path
	}

	return path
}
