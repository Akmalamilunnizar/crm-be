package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AdminGeocodingRepositoryInterface interface {
	ReverseGeocodeRepository(lat, lng string) (map[string]interface{}, error)
}

type AdminGeocodingRepositoryStruct struct {
	// No database needed for this service
}

func NewAdminGeocodingRepository() AdminGeocodingRepositoryInterface {
	return &AdminGeocodingRepositoryStruct{}
}

func (r AdminGeocodingRepositoryStruct) ReverseGeocodeRepository(lat, lng string) (map[string]interface{}, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build Nominatim API URL
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s&addressdetails=1", lat, lng)

	// Create request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers to avoid rate limiting
	req.Header.Set("User-Agent", "CRM-System/1.0")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result, nil
}
