package geocoding

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type AdminGeocodingHandlerInterface interface {
	ReverseGeocodeHandler(c *fiber.Ctx) error
}

type AdminGeocodingHandlerStruct struct {
	service AdminGeocodingServiceInterface
}

func NewAdminGeocodingHandler(service AdminGeocodingServiceInterface) AdminGeocodingHandlerInterface {
	return &AdminGeocodingHandlerStruct{service: service}
}

func (h AdminGeocodingHandlerStruct) ReverseGeocodeHandler(c *fiber.Ctx) error {
	lat := c.Query("lat")
	lng := c.Query("lng")

	if lat == "" || lng == "" {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Latitude and longitude are required", nil)
	}

	result, err := h.service.ReverseGeocodeService(lat, lng)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Reverse geocoding successful", result)
}
