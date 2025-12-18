package helpers

import (
	"fmt"
	"log/slog"
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
	// Try to load .env file, but don't fail if it doesn't exist (for production)
	err := godotenv.Load()
	if err != nil {
		slog.Debug("No .env file found, using environment variables")
	}
	// Add a new global variable for the secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                             // Subject (user identifier)
		"iss": "lillyapps",                          // Issuer
		"aud": role,                                 // Audience (user role)
		"exp": time.Now().Add(2 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                    // Issued at
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
	_ = godotenv.Load()
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token not provided", nil)
	}

	// 1. Parse JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token (parse)", err.Error())
	}

	if !token.Valid {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token (not valid)", nil)
	}

	dbmysql := database.GetDB()
	var user entities.User

	// 2. Cek token di DB
	if err := dbmysql.First(&user, "token = ?", tokenString).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token (db not found)", nil)
		}
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token (db error)", err.Error())
	}

	// 3. Cek exp (juga perbaiki tipe)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if expVal, ok := claims["exp"].(float64); ok {
			exp := int64(expVal)
			if time.Unix(exp, 0).Before(time.Now()) {
				return ResponseUtils(c, fiber.StatusUnauthorized, false, "Token Expired", nil)
			}
		}
	} else {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims", nil)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims (sub)", nil)
	}

	aud, err := token.Claims.GetAudience()
	if err != nil || len(aud) == 0 {
		return ResponseUtils(c, fiber.StatusUnauthorized, false, "Invalid Token Claims (aud)", nil)
	}

	c.Locals("user_id", subject)
	c.Locals("role", aud[0])

	return c.Next()
}

func CustomerVerifyToken(c *fiber.Ctx) error {
	// Try to load .env file, but don't fail if it doesn't exist (for production/Docker)
	err := godotenv.Load()
	if err != nil {
		slog.Debug("No .env file found, using environment variables")
	}
	// Add a new global variable for the secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString := c.Get("Authorization")
	slog.Debug("Token from header", "token", tokenString)
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
