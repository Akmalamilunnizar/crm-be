package wa

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type WAHandler struct {
	service *WAService
}

func NewWAHandler() *WAHandler {
	return &WAHandler{
		service: NewWAService(),
	}
}

type SendMessageRequest struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type SendMessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		To      string `json:"to"`
		Message string `json:"message"`
		SentAt  string `json:"sent_at"`
	} `json:"data"`
}

func (h *WAHandler) SendMessage(c *fiber.Ctx) error {
	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.To == "" || req.Message == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "to and message are required",
		})
	}

	// Send the message using WhatsApp service
	err := h.service.SendMessage(req.To, req.Message)
	if err != nil {
		log.Printf("[WA] Failed to send message to %s: %v", req.To, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to send message: " + err.Error(),
		})
	}

	// Success response
	response := SendMessageResponse{
		Status:  "success",
		Message: "Message sent successfully",
	}
	response.Data.To = req.To
	response.Data.Message = req.Message
	response.Data.SentAt = time.Now().Format(time.RFC3339)

	return c.JSON(response)
}
