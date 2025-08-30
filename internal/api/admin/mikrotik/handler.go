package mikrotik

import (
	"net/http"
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
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusOK, true, "Status retrieved", map[string]interface{}{
			"status":  "disconnected",
			"message": "Not connected to any MikroTik device",
			"host":    "",
			"port":    0,
		})
	}

	status := h.mikroTikService.GetConnectionStatus()
	return helpers.ResponseUtils(c, http.StatusOK, true, "Status retrieved", status)
}

// Get MikroTik logs
func (h *MikroTikHandler) GetLogs(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	timeRange := c.Query("timeRange", "1d")

	logs, err := h.mikroTikService.GetLogs(timeRange)
	if err != nil {
		return helpers.ResponseUtils(c, http.StatusInternalServerError, false, "Failed to get logs: "+err.Error(), nil)
	}

	return helpers.ResponseUtils(c, http.StatusOK, true, "Logs retrieved successfully", map[string]interface{}{
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

// Get real-time logs with WebSocket support
func (h *MikroTikHandler) GetRealTimeLogs(c *fiber.Ctx) error {
	if h.mikroTikService == nil {
		return helpers.ResponseUtils(c, http.StatusBadRequest, false, "Not connected to any MikroTik device", nil)
	}

	// This would typically be handled by WebSocket middleware
	// For now, we'll return the latest logs
	timeRange := c.Query("timeRange", "1h")

	logs, err := h.mikroTikService.GetLogs(timeRange)
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
