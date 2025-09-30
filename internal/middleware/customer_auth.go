package middleware

import (
	"log"
	"os"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// CustomerAuthMiddleware validates JWT tokens for customer access
func CustomerAuthMiddleware(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	log.Println("Customer token from header:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Token not provided", nil)
	}

	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is valid
	if !token.Valid {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(int64); ok {
			if time.Unix(exp, 0).Before(time.Now()) {
				return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Token Expired", nil)
			}
		}
	} else {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	audience, err := token.Claims.GetAudience()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	// For customer authentication, we need to verify the customer exists
	db := database.GetDB()
	var customer entities.Customer
	if err := db.Where("id = ? OR email = ?", subject, subject).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Customer not found", nil)
		}
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, "Database error", nil)
	}

	// Set customer information in context
	c.Locals("customer_id", customer.ID)
	c.Locals("customer_phone", customer.Phone)
	c.Locals("customer_name", customer.Name)
	c.Locals("user_id", subject)
	c.Locals("role", audience[0])

	return c.Next()
}

// AdminAuthMiddleware validates JWT tokens for admin/staff access
func AdminAuthMiddleware(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	log.Println("Admin token from header:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Token not provided", nil)
	}

	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is valid
	if !token.Valid {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(int64); ok {
			if time.Unix(exp, 0).Before(time.Now()) {
				return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Token Expired", nil)
			}
		}
	} else {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	audience, err := token.Claims.GetAudience()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	// For admin authentication, verify the user exists in the users table
	db := database.GetDB()
	var user entities.User
	if err := db.Where("id = ? OR email = ?", subject, subject).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "User not found", nil)
		}
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, "Database error", nil)
	}

	// Set user information in context
	c.Locals("user_id", user.ID)
	c.Locals("user_email", user.Email)
	c.Locals("user_name", user.Name)
	c.Locals("role", audience[0])

	return c.Next()
}
