package helpers

import (
	"fmt"
	"log"
	"os"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Function to create JWT tokens with claims
func CreateToken(username string, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Add a new global variable for the secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "lillyapps",                      // Issuer
		"aud": role,                             // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims, secretKey, username)
	return tokenString, nil
}

// Function to verify JWT tokens
func VerifyToken(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Add a new global variable for the secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	log.Println("Token from header:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token not provided", nil)
	}
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is valid
	if !token.Valid {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)

	}

	dbmysql := database.GetDB()
	// Check if the token exists in the database
	var user entities.User
	if err := dbmysql.First(&user, "token = ?", tokenString).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)

		}
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)

	}
	// Check if the token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(int64); ok {
			if time.Unix(exp, 0).Before(time.Now()) {
				return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token Expired", nil)

			}
		}
	} else {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)

	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	audience, err := token.Claims.GetAudience()
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}
	c.Locals("user_id", subject)
	c.Locals("role", audience[0])
	return c.Next()
}

func CustomerVerifyToken(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Add a new global variable for the secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	log.Println("Token from header:", tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token not provided", nil)
	}
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)
	}

	// Check if the token is valid
	if !token.Valid {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token", nil)

	}

	// Check if the token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(int64); ok {
			if time.Unix(exp, 0).Before(time.Now()) {
				return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token Expired", nil)

			}
		}
	} else {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	audience, err := token.Claims.GetAudience()
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	c.Locals("user_id", subject)
	c.Locals("role", audience[0])
	return c.Next()
}

func VerifyTokenFromDB(tokenString string) error {

	return nil
}
