package telegramapi

import (
	"log"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
	"skripsi-be/internal/telegram"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	telegramService *telegram.Service
}

func NewHandler() *Handler {
	return &Handler{
		telegramService: telegram.NewService(),
	}
}

// TestNotification sends a test notification to verify Telegram bot connectivity
func (h *Handler) TestNotification(c *fiber.Ctx) error {
	// Parse request
	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	// Validate role
	validRoles := []string{"TECHNICIAN", "NOC", "CUSTOMER SERVICE"}
	roleValid := false
	for _, validRole := range validRoles {
		if req.Role == validRole {
			roleValid = true
			break
		}
	}

	if !roleValid {
		return helpers.ResponseUtils(c, 400, false, "Invalid role. Valid roles: TECHNICIAN, NOC, CUSTOMER SERVICE", nil)
	}

	// Send test notification
	err := h.telegramService.TestConnection(req.Role)
	if err != nil {
		log.Printf("Failed to send test notification: %v", err)
		return helpers.ResponseUtils(c, 500, false, "Failed to send test notification: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Test notification sent successfully", map[string]string{
		"role":    req.Role,
		"message": "Test notification sent to " + req.Role + " channel",
	})
}

// TestTicketNotification sends a test ticket notification
func (h *Handler) TestTicketNotification(c *fiber.Ctx) error {
	// Parse request
	var req struct {
		Role   string `json:"role"`
		Action string `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	// Validate role
	validRoles := []string{"TECHNICIAN", "NOC", "CUSTOMER SERVICE"}
	roleValid := false
	for _, validRole := range validRoles {
		if req.Role == validRole {
			roleValid = true
			break
		}
	}

	if !roleValid {
		return helpers.ResponseUtils(c, 400, false, "Invalid role. Valid roles: TECHNICIAN, NOC, CUSTOMER SERVICE", nil)
	}

	// Create a mock ticket for testing
	now := time.Now()
	mockTicket := &entities.TroubleTicket{
		ID:          999,
		CustomerID:  "TEST_CUSTOMER_001",
		Title:       "Test Ticket - Telegram Notification",
		Description: stringPtr("This is a test ticket to verify Telegram notifications are working correctly."),
		Type:        stringPtr("test"),
		Status:      "ongoing",
		CreatedAt:   &now,
	}

	action := req.Action
	if action == "" {
		action = "Test Notification"
	}

	// Send test ticket notification
	err := h.telegramService.SendTicketNotification(req.Role, mockTicket, action)
	if err != nil {
		log.Printf("Failed to send test ticket notification: %v", err)
		return helpers.ResponseUtils(c, 500, false, "Failed to send test ticket notification: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Test ticket notification sent successfully", map[string]interface{}{
		"role":   req.Role,
		"action": action,
		"ticket": mockTicket,
	})
}

// GetChannelInfo returns the channel mapping information
func (h *Handler) GetChannelInfo(c *fiber.Ctx) error {
	channelInfo := map[string]interface{}{
		"channels": map[string]string{
			"TECHNICIAN":       "TEKNISI CRM (-1002731855415)",
			"NOC":              "NOC CRM (-1002927002181)",
			"CUSTOMER SERVICE": "CS CRM (-1002907098486)",
		},
		"bot_configured": h.telegramService != nil,
	}

	return helpers.ResponseUtils(c, 200, true, "Channel information retrieved", channelInfo)
}

// stringPtr returns a pointer to string
func stringPtr(s string) *string {
	return &s
}
