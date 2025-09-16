package monitoring

import (
	"skripsi-be/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}


// GetConnectionHistory returns connection history for the authenticated customer
func (h *Handler) GetConnectionHistory(c *fiber.Ctx) error {
	// Get customer ID from JWT token
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return helpers.ResponseUtils(c, 401, false, "Customer ID not found in token", nil)
	}
	
	customerID, ok := customerIDVal.(string)
	if !ok || customerID == "" {
		return helpers.ResponseUtils(c, 401, false, "Invalid customer ID in token", nil)
	}

	// Get time range from query parameter
	timeRange := c.Query("timeRange", "1d")

	history, err := h.svc.GetConnectionHistory(customerID, timeRange)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "ok", history)
}


// Response structures for API responses
type ConnectionEvent struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration,omitempty"`
	Message   string    `json:"message"`
}
