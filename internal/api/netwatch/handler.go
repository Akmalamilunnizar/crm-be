package netwatchapi

import (
	"skripsi-be/internal/api/netwatch/service"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s}
}

// ListDevices returns all Netwatch devices
func (h *Handler) ListDevices(c *fiber.Ctx) error {
	devices, err := h.svc.ListDevices()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", devices)
}

// GetDevice returns a specific device by ID
func (h *Handler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	device, err := h.svc.GetDevice(id)
	if err != nil {
		return helpers.ResponseUtils(c, 404, false, "Device not found", nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", device)
}

// CreateDevice creates a new Netwatch device
func (h *Handler) CreateDevice(c *fiber.Ctx) error {
	var device entities.NetwatchDevice
	if err := c.BodyParser(&device); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	createdDevice, err := h.svc.CreateDevice(&device)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 201, true, "Device created", createdDevice)
}

// UpdateDevice updates an existing device
func (h *Handler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device entities.NetwatchDevice
	if err := c.BodyParser(&device); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	device.ID = id
	updatedDevice, err := h.svc.UpdateDevice(&device)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Device updated", updatedDevice)
}

// DeleteDevice deletes a device
func (h *Handler) DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteDevice(id); err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Device deleted", nil)
}

// ListEvents returns all Netwatch events
func (h *Handler) ListEvents(c *fiber.Ctx) error {
	events, err := h.svc.ListEvents()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", events)
}

// GetEvent returns a specific event by ID
func (h *Handler) GetEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := h.svc.GetEvent(id)
	if err != nil {
		return helpers.ResponseUtils(c, 404, false, "Event not found", nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", event)
}

// CreateEvent creates a new Netwatch event (for webhook/syslog)
func (h *Handler) CreateEvent(c *fiber.Ctx) error {
	var event entities.NetwatchEvent
	if err := c.BodyParser(&event); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	createdEvent, err := h.svc.CreateEvent(&event)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 201, true, "Event created", createdEvent)
}

// SyncDevices syncs devices from MikroTik
func (h *Handler) SyncDevices(c *fiber.Ctx) error {
	configID := c.Query("config_id")
	if configID == "" {
		return helpers.ResponseUtils(c, 400, false, "config_id is required", nil)
	}

	if err := h.svc.SyncDevicesFromMikroTik(configID); err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Devices synced successfully", nil)
}

// ListConfigs returns all Netwatch configurations
func (h *Handler) ListConfigs(c *fiber.Ctx) error {
	configs, err := h.svc.ListConfigs()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", configs)
}

// CreateConfig creates a new Netwatch configuration
func (h *Handler) CreateConfig(c *fiber.Ctx) error {
	var config entities.NetwatchConfig
	if err := c.BodyParser(&config); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	createdConfig, err := h.svc.CreateConfig(&config)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 201, true, "Configuration created", createdConfig)
}

// UpdateConfig updates an existing configuration
func (h *Handler) UpdateConfig(c *fiber.Ctx) error {
	id := c.Params("id")
	var config entities.NetwatchConfig
	if err := c.BodyParser(&config); err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid request body", nil)
	}

	config.ID = id
	updatedConfig, err := h.svc.UpdateConfig(&config)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Configuration updated", updatedConfig)
}

// DeleteConfig deletes a configuration
func (h *Handler) DeleteConfig(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteConfig(id); err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Configuration deleted", nil)
}

// GetMonitoringStatus returns the current monitoring status
func (h *Handler) GetMonitoringStatus(c *fiber.Ctx) error {
	status, err := h.svc.GetMonitoringStatus()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", status)
}
