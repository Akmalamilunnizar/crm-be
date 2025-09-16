package main

import (
	"encoding/json"
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	database.GetDB()
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		BodyLimit:    50 * 1024 * 1024,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	// Initialize default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Recovery and request logging for diagnostics
	app.Use(recover.New())
	app.Use(logger.New())

	// Initialize routes
	routes.RouteFiber(app)

	log.Fatal(app.Listen(":3001"))
}
