package wa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	return &WAService{
		apiURL:   os.Getenv("WHATSAPP_API_URL"),
		apiToken: os.Getenv("WHATSAPP_API_TOKEN"),
		phoneID:  os.Getenv("WHATSAPP_PHONE_ID"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *WAService) SendMessage(to, message string) error {
	// If no API configuration, just log (mock mode)
	if s.apiURL == "" || s.apiToken == "" || s.phoneID == "" {
		fmt.Printf("[WA MOCK] Would send to %s: %s\n", to, message)
		fmt.Printf("[WA MOCK] To enable real messaging, set environment variables:\n")
		fmt.Printf("[WA MOCK] - WHATSAPP_API_URL=https://graph.facebook.com/v18.0\n")
		fmt.Printf("[WA MOCK] - WHATSAPP_API_TOKEN=your_meta_access_token\n")
		fmt.Printf("[WA MOCK] - WHATSAPP_PHONE_ID=your_phone_number_id\n")
		return nil
	}

	// Ensure phone number format (add country code if missing)
	if !s.isValidPhoneNumber(to) {
		return fmt.Errorf("invalid phone number format: %s (use format: 628123456789)", to)
	}

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

	// Verbose request logging (safe: token masked)
	maskedToken := ""
	if len(s.apiToken) > 12 {
		maskedToken = s.apiToken[:8] + "***" + s.apiToken[len(s.apiToken)-4:]
	}
	fmt.Printf("[WA DEBUG] POST %s\n", url)
	fmt.Printf("[WA DEBUG] PhoneID=%s To=%s Type=text\n", s.phoneID, to)
	fmt.Printf("[WA DEBUG] Authorization=Bearer %s\n", maskedToken)
	fmt.Printf("[WA DEBUG] Payload=%s\n", string(jsonData))
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

	// Verbose response logging
	fmt.Printf("[WA DEBUG] Status=%d\n", resp.StatusCode)
	fmt.Printf("[WA DEBUG] Response=%s\n", string(body))

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Parse error response for better error messages
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
			} `json:"error"`
		}

		if err := json.Unmarshal(body, &errorResp); err == nil {
			// Handle specific error codes
			// switch errorResp.Error.Code {
			// case 131030:
			// 	return fmt.Errorf("phone number %s is not in the allowed list. Please add it to your Meta WhatsApp Business API allowed recipients list", to)
			// case 190:
			// 	return fmt.Errorf("access token expired. Please update your WHATSAPP_API_TOKEN")
			// case 368:
			// 	return fmt.Errorf("rate limit exceeded. Please wait before sending more messages")
			// default:
			// 	return fmt.Errorf("WhatsApp API error: %s", errorResp.Error.Message)
			// }
		}

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

	// 200 but no messages array – log as anomaly for visibility
	fmt.Printf("[WA WARN] 200 OK but no messages id returned. Contacts=%v\n", apiResp.Contacts)
	return fmt.Errorf("Meta responded 200 but no message id returned for %s", to)
}

// isValidPhoneNumber checks if the phone number is in the correct format
func (s *WAService) isValidPhoneNumber(phone string) bool {
	// Basic validation: should start with country code and be numeric
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	// Check if it starts with a country code (like 62 for Indonesia)
	if phone[0] != '+' && (len(phone) < 10 || len(phone) > 15) {
		return false
	}

	return true
}
