package geocoding

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func AdminGeocodingRoute(app fiber.Router) {
	repository := NewAdminGeocodingRepository()
	service := NewAdminGeocodingService(repository)
	handler := NewAdminGeocodingHandler(service)

	// Add authentication middleware
	app.Use(helpers.VerifyToken)

	// Geocoding routes
	app.Get("reverse-geocode", handler.ReverseGeocodeHandler)
}
