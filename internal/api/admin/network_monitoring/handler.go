package network_monitoring

import (
	"net/http"

	"skripsi-be/internal/helpers"
	"skripsi-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

type NetworkMonitoringHandler struct {
	networkMonitoringService *services.NetworkMonitoringAssistant
}

func NewNetworkMonitoringHandler(networkMonitoringService *services.NetworkMonitoringAssistant) *NetworkMonitoringHandler {
	return &NetworkMonitoringHandler{
		networkMonitoringService: networkMonitoringService,
	}
}

// MonitorCustomerDevice monitors a specific customer device
func (h *NetworkMonitoringHandler) MonitorCustomerDevice(c *fiber.Ctx) error {
	customerID := c.Params("customerId")
	ipAddress := c.Query("ip")

	if customerID == "" || ipAddress == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Customer ID and IP address are required", nil)
	}

	status, err := h.networkMonitoringService.MonitorCustomerDevice(customerID, ipAddress)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to monitor device: "+err.Error(), nil)
	}

	// Format response for frontend
	formattedResponse := h.networkMonitoringService.FormatStatusResponse(status)

	return helpers.ResponseUtils(c, http.StatusOK, true, "Device status checked successfully", map[string]interface{}{
		"status":      status.Status,
		"ip":          status.IP,
		"last_checked": status.LastChecked,
		"recommendation": status.Recommendation,
		"formatted_response": formattedResponse,
	})
}

// GetCustomerConnectionStatus returns connection status for a customer
func (h *NetworkMonitoringHandler) GetCustomerConnectionStatus(c *fiber.Ctx) error {
	customerID := c.Params("customerId")

	if customerID == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Customer ID is required", nil)
	}

	statuses, err := h.networkMonitoringService.GetCustomerConnectionStatus(customerID)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get connection status: "+err.Error(), nil)
	}

	// Format responses for frontend
	var formattedResponses []string
	for _, status := range statuses {
		formattedResponses = append(formattedResponses, h.networkMonitoringService.FormatStatusResponse(&status))
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Connection status retrieved successfully", map[string]interface{}{
		"devices": statuses,
		"formatted_responses": formattedResponses,
	})
}

// GenerateNetwatchScript generates a MikroTik Netwatch script for an IP
func (h *NetworkMonitoringHandler) GenerateNetwatchScript(c *fiber.Ctx) error {
	var req struct {
		IPAddress    string `json:"ip_address"`
		CustomerName string `json:"customer_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	if req.IPAddress == "" || req.CustomerName == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "IP address and customer name are required", nil)
	}

	script := h.networkMonitoringService.GenerateNetwatchScript(req.IPAddress, req.CustomerName)

	return helpers.ResponseUtils(c, http.StatusOK, true, "Netwatch script generated successfully", map[string]interface{}{
		"script": script,
		"ip_address": req.IPAddress,
		"customer_name": req.CustomerName,
	})
}

// GetTroubleshootingRecommendations returns troubleshooting tips for a down device
func (h *NetworkMonitoringHandler) GetTroubleshootingRecommendations(c *fiber.Ctx) error {
	ipAddress := c.Query("ip")

	if ipAddress == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "IP address is required", nil)
	}

	// Generate recommendations
	recommendations := []string{
		"Periksa koneksi kabel LAN/WiFi",
		"Restart router/modem pelanggan",
		"Periksa status listrik di lokasi pelanggan",
		"Verifikasi konfigurasi IP address",
		"Periksa firewall atau security settings",
		"Hubungi teknisi untuk pengecekan fisik",
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Troubleshooting recommendations retrieved", map[string]interface{}{
		"ip_address": ipAddress,
		"recommendations": recommendations,
		"formatted_response": "[DOWN] " + ipAddress + " – device offline, rekomendasi: " + recommendations[0] + ", " + recommendations[1] + ", " + recommendations[2],
	})
}

// BulkMonitorDevices monitors multiple devices at once
func (h *NetworkMonitoringHandler) BulkMonitorDevices(c *fiber.Ctx) error {
	var req struct {
		Devices []struct {
			CustomerID string `json:"customer_id"`
			IPAddress  string `json:"ip_address"`
		} `json:"devices"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	if len(req.Devices) == 0 {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "At least one device is required", nil)
	}

	var results []map[string]interface{}
	var formattedResponses []string

	for _, device := range req.Devices {
		status, err := h.networkMonitoringService.MonitorCustomerDevice(device.CustomerID, device.IPAddress)
		if err != nil {
			results = append(results, map[string]interface{}{
				"customer_id": device.CustomerID,
				"ip_address":  device.IPAddress,
				"status":      "error",
				"error":       err.Error(),
			})
			formattedResponses = append(formattedResponses, "[ERROR] "+device.IPAddress+" – monitoring failed: "+err.Error())
			continue
		}

		formattedResponse := h.networkMonitoringService.FormatStatusResponse(status)
		formattedResponses = append(formattedResponses, formattedResponse)

		results = append(results, map[string]interface{}{
			"customer_id":    device.CustomerID,
			"ip_address":     status.IP,
			"status":         status.Status,
			"last_checked":   status.LastChecked,
			"recommendation": status.Recommendation,
		})
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Bulk monitoring completed", map[string]interface{}{
		"results":            results,
		"formatted_responses": formattedResponses,
		"total_devices":      len(req.Devices),
		"successful_checks":  len(results),
	})
}

// GetMonitoringStats returns statistics about monitored devices
func (h *NetworkMonitoringHandler) GetMonitoringStats(c *fiber.Ctx) error {
	// This would typically query the database for statistics
	// For now, return basic stats structure
	stats := map[string]interface{}{
		"total_monitored_devices": 0,
		"devices_up":              0,
		"devices_down":            0,
		"last_check_time":         nil,
		"monitoring_active":       true,
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Monitoring statistics retrieved", stats)
}
