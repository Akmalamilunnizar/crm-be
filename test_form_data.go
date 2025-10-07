package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Test endpoint to debug form data
	app.Post("/test-form", func(c *fiber.Ctx) error {
		log.Printf("=== TEST FORM DATA DEBUG ===")

		// Parse multipart form
		form, err := c.MultipartForm()
		if err != nil {
			log.Printf("Error parsing multipart form: %v", err)
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid multipart form data",
				"details": err.Error(),
			})
		}

		// Log all form values
		log.Printf("Form values:")
		for key, values := range form.Value {
			log.Printf("  %s: %v", key, values)
		}

		// Log specific fields we care about
		macAddress := c.FormValue("mac_address")
		assetItemID := c.FormValue("asset_item_id")
		assetsID := c.FormValue("assets_id")

		log.Printf("=== SPECIFIC FIELDS ===")
		log.Printf("mac_address: '%s'", macAddress)
		log.Printf("asset_item_id: '%s'", assetItemID)
		log.Printf("assets_id: '%s'", assetsID)
		log.Printf("=== END DEBUG ===")

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"mac_address":   macAddress,
				"asset_item_id": assetItemID,
				"assets_id":     assetsID,
				"all_fields":    form.Value,
			},
		})
	})

	log.Println("Test server starting on :3002")
	log.Fatal(app.Listen(":3002"))
}
