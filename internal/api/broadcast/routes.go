package broadcast

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupBroadcastRoutes sets up all broadcast-related routes
func SetupBroadcastRoutes(app *fiber.App, db *gorm.DB) {
	// Create broadcast route group
	broadcast := app.Group("/api/broadcast")

	// Add authentication middleware
	broadcast.Use(helpers.VerifyToken)

	// Register routes
	broadcast.Post("/send", func(c *fiber.Ctx) error {
		return HandleSendBroadcast(c, db)
	})

	broadcast.Get("/history", func(c *fiber.Ctx) error {
		return HandleGetBroadcastHistory(c, db)
	})

	// Get recipients from a specific broadcast (for follow-up broadcasts)
	broadcast.Get("/:id/recipients", func(c *fiber.Ctx) error {
		return HandleGetBroadcastRecipients(c, db)
	})
}
