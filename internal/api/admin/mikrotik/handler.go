package mikrotik

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"skripsi-be/internal/helpers"
	"skripsi-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MikroTikHandler struct {
	mikroTikService *services.MikroTikService
}

func NewMikroTikHandler(mikroTikService *services.MikroTikService) *MikroTikHandler {
	return &MikroTikHandler{
		mikroTikService: mikroTikService,
	}
}

// Connect to MikroTik device
func (h *MikroTikHandler) Connect(c *fiber.Ctx) error {
	var req struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	// Validate required fields
	if req.Host == "" || req.Port == 0 || req.Username == "" || req.Password == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Host, port, username, and password are required", nil)
	}

	// Create new MikroTik service with provided config
	config := &services.MikroTikConfig{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}

	h.mikroTikService = services.NewMikroTikService(config)

	// Attempt to connect
	if err := h.mikroTikService.Connect(); err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to connect to MikroTik: "+err.Error(), nil)
	}

	// expose shared instance for background jobs
	services.SetSharedMikroTikService(h.mikroTikService)

	return helpers.ResponseUtils(c, http.StatusOK, true, "Successfully connected to MikroTik", map[string]interface{}{
		"host":   req.Host,
		"port":   req.Port,
		"status": "connected",
	})
}

// Disconnect from MikroTik device
func (h *MikroTikHandler) Disconnect(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	if err := h.mikroTikService.Disconnect(); err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to disconnect: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Successfully disconnected from MikroTik", nil)
}

// Get connection status
func (h *MikroTikHandler) GetStatus(c *fiber.Ctx) error {
	// Always check the shared service (set by auto-connect)
	sharedService := services.GetSharedMikroTikService()
	if sharedService == nil {
		return helpers.ResponseUtils(c, http.StatusOK, true, "Status retrieved", map[string]interface{}{
			"status":  "disconnected",
			"message": "Not connected to any MikroTik device",
			"host":    "",
			"port":    0,
		})
	}

	status := sharedService.GetConnectionStatus()
	return helpers.ResponseUtils(c, http.StatusOK, true, "Status retrieved", status)
}

// Get MikroTik logs (Netwatch logs)
func (h *MikroTikHandler) GetLogs(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	timeRange := c.Query("timeRange", "1d")

	// Use Netwatch logs instead of general logs
	logs, err := h.mikroTikService.GetNetwatchLogs(timeRange)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get netwatch logs: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Netwatch logs retrieved successfully", map[string]interface{}{
		"logs":      logs,
		"count":     len(logs),
		"timeRange": timeRange,
	})
}

// Execute custom command on MikroTik
func (h *MikroTikHandler) ExecuteCommand(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	var req struct {
		Command string `json:"command"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	if req.Command == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Command is required", nil)
	}

	output, err := h.mikroTikService.ExecuteCommand(req.Command)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to execute command: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Command executed successfully", map[string]interface{}{
		"command": req.Command,
		"output":  output,
	})
}

// Get system information
func (h *MikroTikHandler) GetSystemInfo(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	info, err := h.mikroTikService.GetSystemInfo()
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get system info: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "System information retrieved successfully", info)
}

// Set hotspot IP binding type (for customer isolation)
func (h *MikroTikHandler) SetHotspotIPBindingType(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	var req struct {
		MacAddress string `json:"mac_address" validate:"required"`
		Type       string `json:"type" validate:"required,oneof=regular bypassed"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	if req.MacAddress == "" || req.Type == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "MAC address and type are required", nil)
	}

	err := h.mikroTikService.SetHotspotIPBindingType(req.MacAddress, req.Type)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to set IP binding type: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "IP binding type updated successfully", map[string]interface{}{
		"mac_address": req.MacAddress,
		"type":        req.Type,
	})
}

// Get hotspot IP bindings
func (h *MikroTikHandler) GetHotspotIPBindings(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	bindings, err := h.mikroTikService.GetHotspotIPBindings()
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get IP bindings: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "IP bindings retrieved successfully", bindings)
}

// Get specific hotspot IP binding by MAC address
func (h *MikroTikHandler) GetHotspotIPBindingByMAC(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	macAddress := c.Params("mac")
	if macAddress == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "MAC address is required", nil)
	}

	binding, err := h.mikroTikService.GetHotspotIPBindingByMAC(macAddress)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusNotFound, false, "IP binding not found: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "IP binding retrieved successfully", binding)
}

// Get Netwatch devices
func (h *MikroTikHandler) GetNetwatchDevices(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	devices, err := h.mikroTikService.GetNetwatchDevices()
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get netwatch devices: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Netwatch devices retrieved successfully", map[string]interface{}{
		"devices": devices,
		"count":   len(devices),
	})
}

// Get real-time logs with WebSocket support
func (h *MikroTikHandler) GetRealTimeLogs(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	// This would typically be handled by WebSocket middleware
	// For now, we'll return the latest logs
	timeRange := c.Query("timeRange", "1h")

	logs, err := h.mikroTikService.GetNetwatchLogs(timeRange)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get logs: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Real-time logs retrieved successfully", map[string]interface{}{
		"logs":      logs,
		"count":     len(logs),
		"timeRange": timeRange,
		"timestamp": time.Now(),
	})
}

// GetDHCPLease fetches DHCP lease information for a given MAC address
func (h *MikroTikHandler) GetDHCPLease(c *fiber.Ctx) error {
	var req struct {
		MacAddress string `json:"mac_address"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Invalid request body", nil)
	}

	if req.MacAddress == "" {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "MAC address is required", nil)
	}

	// Always use the shared MikroTik service (which is set by auto-connect)
	sharedService := services.GetSharedMikroTikService()
	if sharedService == nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "MikroTik service not initialized. Please connect to MikroTik first.", nil)
	}

	// Check if MikroTik service is connected
	if !sharedService.IsConnected() {
		return helpers.ResponseUtils(c, http.StatusServiceUnavailable, false, "MikroTik service not connected. Please connect to MikroTik first.", nil)
	}

	// Execute MikroTik command to get DHCP lease
	cmd := "/ip/dhcp-server/lease/print where mac-address=" + req.MacAddress + " status=bound disabled=no"
	output, err := sharedService.ExecuteCommand(cmd)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to query DHCP lease: "+err.Error(), map[string]interface{}{
			"mac_address": req.MacAddress,
			"command":     cmd,
			"error":       err.Error(),
		})
	}

	// Parse the output to extract IP address
	ipAddress := parseDHCPLeaseOutput(output)
	if ipAddress == "" {
		return helpers.ResponseUtils(c, http.StatusNotFound, false, "No active DHCP lease found for MAC address "+req.MacAddress, map[string]interface{}{
			"mac_address":   req.MacAddress,
			"command":       cmd,
			"raw_output":    output,
			"parsing_debug": "Failed to extract IP address from output",
		})
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "DHCP lease found successfully", map[string]interface{}{
		"mac_address": req.MacAddress,
		"ip_address":  ipAddress,
		"command":     cmd,
		"raw_output":  output,
	})
}

// parseDHCPLeaseOutput parses MikroTik output to extract IP address
func parseDHCPLeaseOutput(output string) string {
	// MikroTik output format can be different:
	// Format 1: "address=10.10.21.125 mac-address=40:EE:15:7D:43:01 client-id=1:40:ee:15:7d:43:1 server=dhcp1 status=bound"
	// Format 2: Table format with columns like "10.10.21.125 40:EE:15:7D:43:01 dhcp1 bound 7m2s"

	lines := strings.Split(output, "\n")

	// Debug: log the output for troubleshooting
	fmt.Printf("DEBUG: Parsing DHCP output (%d lines):\n", len(lines))
	for i, line := range lines {
		fmt.Printf("  Line %d: '%s'\n", i, line)
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip header lines
		if strings.Contains(line, "Columns:") || strings.Contains(line, "# ADDRESS") || strings.Contains(line, "Flags:") {
			continue
		}

		// Try format 1: address=10.10.21.125 mac-address=...
		re1 := regexp.MustCompile(`address=([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
		matches1 := re1.FindStringSubmatch(line)
		if len(matches1) > 1 {
			return matches1[1]
		}

		// Try format 2: Table format - IP is usually the 3rd field (after number and flags)
		// Look for lines that start with a number (skip # entries)
		if !strings.HasPrefix(line, "#") && !strings.Contains(line, "Flags:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				// Check if first field is a number and second field is flags (like "D")
				// Then the third field should be the IP address
				if len(fields[0]) > 0 && len(fields[1]) == 1 && len(fields[2]) > 0 {
					ipField := fields[2]
					// Validate it's an IP address
					ipRegex := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
					if ipRegex.MatchString(ipField) {
						fmt.Printf("DEBUG: Found IP address: %s\n", ipField)
						return ipField
					}
				}

				// Fallback: check all fields for IP addresses
				for _, field := range fields {
					ipRegex := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
					if ipRegex.MatchString(field) {
						fmt.Printf("DEBUG: Found IP address in field: %s\n", field)
						return field
					}
				}
			}
		}
	}

	return ""
}
