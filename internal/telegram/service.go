package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"skripsi-be/internal/models/entities"
)

type Service struct {
	botToken   string
	channelMap map[string]string
}

type Message struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewService() *Service {
	// Get bot token from environment or use default
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		botToken = "8497644710:AAHYvmT7lr-nfJtwwF3FnuoNadxrXw2361E"
	}

	// Role-based channel mapping - can be overridden by environment variables
	channelMap := map[string]string{
		"TECHNICIAN":       getEnvOrDefault("TELEGRAM_TECHNICIAN_CHANNEL", "-1002731855415"), // TEKNISI CRM
		"NOC":              getEnvOrDefault("TELEGRAM_NOC_CHANNEL", "-1002927002181"),        // NOC CRM
		"CUSTOMER SERVICE": getEnvOrDefault("TELEGRAM_CS_CHANNEL", "-1002907098486"),         // CS CRM
		"CUSTOMER_SERVICE": getEnvOrDefault("TELEGRAM_CS_CHANNEL", "-1002907098486"),         // Alternative CS format
	}

	return &Service{
		botToken:   botToken,
		channelMap: channelMap,
	}
}

// SendTicketNotification sends a notification to the appropriate Telegram channel based on role
func (ts *Service) SendTicketNotification(role string, ticket *entities.TroubleTicket, action string) error {
	// Check if notifications are enabled
	if !ts.isNotificationEnabled() {
		log.Printf("Telegram notifications are disabled")
		return nil
	}

	channelID, exists := ts.channelMap[role]
	if !exists {
		log.Printf("No Telegram channel found for role: %s", role)
		return fmt.Errorf("no telegram channel configured for role: %s", role)
	}

	message := ts.formatTicketMessage(ticket, action)
	return ts.sendMessage(channelID, message)
}

// SendCustomNotification sends a custom message to a specific role's channel
func (ts *Service) SendCustomNotification(role string, message string) error {
	channelID, exists := ts.channelMap[role]
	if !exists {
		log.Printf("No Telegram channel found for role: %s", role)
		return fmt.Errorf("no telegram channel configured for role: %s", role)
	}

	return ts.sendMessage(channelID, message)
}

// formatTicketMessage creates a formatted message for ticket notifications
func (ts *Service) formatTicketMessage(ticket *entities.TroubleTicket, action string) string {
	var message bytes.Buffer

	message.WriteString("🎫 <b>TROUBLE TICKET NOTIFICATION</b>\n\n")
	message.WriteString(fmt.Sprintf("📋 <b>Action:</b> %s\n", action))
	message.WriteString(fmt.Sprintf("🆔 <b>Ticket ID:</b> #%d\n", ticket.ID))
	message.WriteString(fmt.Sprintf("👤 <b>Customer ID:</b> %s\n", ticket.CustomerID))
	message.WriteString(fmt.Sprintf("📝 <b>Title:</b> %s\n", ticket.Title))

	if ticket.Description != nil && *ticket.Description != "" {
		message.WriteString(fmt.Sprintf("📄 <b>Description:</b> %s\n", *ticket.Description))
	}

	if ticket.Type != nil && *ticket.Type != "" {
		message.WriteString(fmt.Sprintf("🔧 <b>Type:</b> %s\n", *ticket.Type))
	}

	message.WriteString(fmt.Sprintf("📊 <b>Status:</b> %s\n", ticket.Status))

	if ticket.CustomerNote != nil && *ticket.CustomerNote != "" {
		message.WriteString(fmt.Sprintf("💬 <b>Customer Note:</b> %s\n", *ticket.CustomerNote))
	}

	if ticket.NOCNote != nil && *ticket.NOCNote != "" {
		message.WriteString(fmt.Sprintf("🔧 <b>NOC Note:</b> %s\n", *ticket.NOCNote))
	}

	if ticket.TechnicianNote != nil && *ticket.TechnicianNote != "" {
		message.WriteString(fmt.Sprintf("🛠️ <b>Technician Note:</b> %s\n", *ticket.TechnicianNote))
	}

	if ticket.CreatedAt != nil {
		message.WriteString(fmt.Sprintf("⏰ <b>Created:</b> %s\n", ticket.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	return message.String()
}

// sendMessage sends a message to a specific Telegram channel
func (ts *Service) sendMessage(channelID string, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", ts.botToken)

	telegramMsg := Message{
		ChatID:    channelID,
		Text:      message,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		log.Printf("Error marshaling Telegram message: %v", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned non-200 status: %d", resp.StatusCode)
		return fmt.Errorf("telegram API error: status %d", resp.StatusCode)
	}

	log.Printf("Successfully sent Telegram notification to channel %s", channelID)
	return nil
}

// GetChannelForRole returns the channel ID for a given role
func (ts *Service) GetChannelForRole(role string) (string, bool) {
	channelID, exists := ts.channelMap[role]
	return channelID, exists
}

// TestConnection sends a test message to verify bot connectivity
func (ts *Service) TestConnection(role string) error {
	testMessage := "🧪 <b>TEST NOTIFICATION</b>\n\nTelegram bot is working correctly for CRM ticket notifications!"
	return ts.SendCustomNotification(role, testMessage)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isNotificationEnabled checks if Telegram notifications are enabled via environment variable
func (ts *Service) isNotificationEnabled() bool {
	enabled := getEnvOrDefault("TELEGRAM_NOTIFICATIONS_ENABLED", "true")
	return enabled == "true" || enabled == "1" || enabled == "yes"
}
