package wa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WAService struct {
	apiURL   string
	apiToken string
	phoneID  string
	client   *http.Client
}

type WhatsAppMessage struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type WhatsAppAPIRequest struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

type WhatsAppAPIResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

func NewWAService() *WAService {
	// Use local whatsapp-web.js service if available, otherwise fall back to Meta API
	serviceURL := os.Getenv("WHATSAPP_SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:3002" // Default whatsapp-web.js service URL
	}

	return &WAService{
		apiURL:   serviceURL,
		apiToken: os.Getenv("WHATSAPP_API_TOKEN"), // Not used for local service
		phoneID:  os.Getenv("WHATSAPP_PHONE_ID"),  // Not used for local service
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// normalizePhoneNumber normalizes Indonesian phone numbers to WhatsApp format
// Converts: 08970833338 -> 628970833338
// Converts: +628970833338 -> 628970833338
// Converts: 628970833338 -> 628970833338 (already correct)
func normalizePhoneNumber(phone string) string {
	// Remove all whitespace and non-digit characters except +
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Remove + if present
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	// If starts with 0, replace with 62 (Indonesia country code)
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	// If doesn't start with 62, add it (assuming Indonesian numbers)
	if !strings.HasPrefix(phone, "62") && len(phone) >= 9 {
		phone = "62" + phone
	}

	return phone
}

func (s *WAService) SendMessage(to, message string) error {
	// Normalize phone number format (remove +, convert 0 to 62)
	normalizedTo := normalizePhoneNumber(to)

	// Ensure phone number format (add country code if missing)
	if !s.isValidPhoneNumber(normalizedTo) {
		return fmt.Errorf("invalid phone number format: %s (normalized: %s, use format: 628123456789)", to, normalizedTo)
	}

	// Check if using local whatsapp-web.js service or Meta API
	// If apiURL doesn't contain "graph.facebook.com", assume it's local service
	usingLocalService := !strings.Contains(s.apiURL, "graph.facebook.com")

	if usingLocalService {
		// Use local whatsapp-web.js service
		return s.sendViaLocalService(normalizedTo, message)
	}

	// Fall back to Meta WhatsApp Business API (if configured)
	if s.apiToken == "" || s.phoneID == "" {
		return fmt.Errorf("WhatsApp service not configured. Please set WHATSAPP_SERVICE_URL or configure Meta API")
	}

	return s.sendViaMetaAPI(normalizedTo, message)
}

// sendViaLocalService sends message via local whatsapp-web.js service
func (s *WAService) sendViaLocalService(to, message string) error {
	// Prepare request for local service
	reqBody := map[string]string{
		"to":      to,
		"message": message,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request to local service
	url := fmt.Sprintf("%s/send-message", s.apiURL)
	fmt.Printf("[WA] Sending via local service to %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to local service: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("WhatsApp service error: %s", errorResp.Message)
		}
		return fmt.Errorf("WhatsApp service error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse success response
	var successResp struct {
		Status string `json:"status"`
		Data   struct {
			MessageID string `json:"message_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &successResp); err == nil {
		fmt.Printf("[WA] Message sent successfully to %s (Message ID: %s)\n", to, successResp.Data.MessageID)
		return nil
	}

	fmt.Printf("[WA] Message sent successfully to %s\n", to)
	return nil
}

// sendViaMetaAPI sends message via Meta WhatsApp Business API (legacy)
func (s *WAService) sendViaMetaAPI(to, message string) error {
	// Prepare the request for Meta WhatsApp Business API
	reqBody := WhatsAppAPIRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
	}
	reqBody.Text.Body = message

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request to Meta Graph API
	url := fmt.Sprintf("%s/%s/messages", s.apiURL, s.phoneID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers for Meta API
	req.Header.Set("Authorization", "Bearer "+s.apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Meta WhatsApp API error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp WhatsAppAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if len(apiResp.Messages) > 0 {
		fmt.Printf("[WA] Message sent successfully to %s (Message ID: %s)\n", to, apiResp.Messages[0].ID)
		return nil
	}

	return fmt.Errorf("Meta responded 200 but no message id returned for %s", to)
}

// isValidPhoneNumber checks if the phone number is in the correct format
// WhatsApp Business API requires: 628123456789 (no +, no leading 0, with country code)
func (s *WAService) isValidPhoneNumber(phone string) bool {
	// Remove + if present for validation
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	// Basic validation: should be 10-15 digits
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	// Check if it's all numeric (after removing +)
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Should start with country code (62 for Indonesia, or other country codes)
	// Accept 62 (Indonesia), 1 (US/Canada), etc.
	if !strings.HasPrefix(phone, "62") && !strings.HasPrefix(phone, "1") {
		// Allow other country codes but must be at least 2 digits
		if len(phone) < 11 {
			return false
		}
	}

	return true
}
