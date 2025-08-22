package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database connection
	db := database.GetDB()

	// Create a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	// Create a test role first (if it doesn't exist)
	testRole := entities.Role{
		ID:   uuid.New().String(),
		Name: "ADMIN",
	}

	// Check if role exists, if not create it
	var existingRole entities.Role
	if err := db.Where("name = ?", "ADMIN").First(&existingRole).Error; err != nil {
		// Role doesn't exist, create it
		if err := db.Create(&testRole).Error; err != nil {
			log.Printf("Error creating role: %v", err)
		} else {
			log.Printf("Created role: %s", testRole.ID)
		}
	} else {
		testRole = existingRole
		log.Printf("Using existing role: %s", testRole.ID)
	}

	// Create test user
	testUser := entities.User{
		ID:       uuid.New().String(),
		Email:    "test@example.com",
		Name:     "Test User",
		Password: string(hashedPassword),
		RoleId:   testRole.ID,
		Phone:    "123456789",
	}

	// Check if user exists
	var existingUser entities.User
	if err := db.Where("email = ?", "test@example.com").First(&existingUser).Error; err != nil {
		// User doesn't exist, create it
		if err := db.Create(&testUser).Error; err != nil {
			log.Fatal("Error creating user:", err)
		}
		log.Printf("Created test user: %s", testUser.Email)
		log.Printf("Password: password123")
	} else {
		log.Printf("User already exists: %s", existingUser.Email)
	}

	log.Println("Test user creation completed!")
}
