package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
		User  struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"user"`
	} `json:"data"`
}

type VerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"user"`
	} `json:"data"`
}

func main() {
	app := fiber.New()

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Auth routes
	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		// Mock login - accept any email/password
		response := LoginResponse{
			Status:  "success",
			Message: "Login successful",
		}
		response.Data.Token = "mock-jwt-token-" + time.Now().Format("20060102150405")
		response.Data.User.ID = "1"
		response.Data.User.Email = req.Email
		response.Data.User.Name = "Mock User"

		return c.JSON(response)
	})

	app.Post("/api/auth/login-customer", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		// Mock customer login
		response := LoginResponse{
			Status:  "success",
			Message: "Customer login successful",
		}
		response.Data.Token = "mock-customer-jwt-token-" + time.Now().Format("20060102150405")
		response.Data.User.ID = "2"
		response.Data.User.Email = req.Email
		response.Data.User.Name = "Mock Customer"

		return c.JSON(response)
	})

	app.Post("/api/auth/verify", func(c *fiber.Ctx) error {
		// Mock token verification
		response := VerifyResponse{
			Status:  "success",
			Message: "Token verified",
		}
		response.Data.User.ID = "1"
		response.Data.User.Email = "admin@example.com"
		response.Data.User.Name = "Mock Admin User"

		return c.JSON(response)
	})

	app.Post("/api/auth/verify-customer", func(c *fiber.Ctx) error {
		// Mock customer token verification
		response := VerifyResponse{
			Status:  "success",
			Message: "Customer token verified",
		}
		response.Data.User.ID = "2"
		response.Data.User.Email = "customer@example.com"
		response.Data.User.Name = "Mock Customer User"

		return c.JSON(response)
	})

	log.Println("Mock server starting on :3001")
	log.Fatal(app.Listen(":3001"))
}
