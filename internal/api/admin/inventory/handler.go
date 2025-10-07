package inventory

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type InventoryHandlerInterface interface {
	CreatePurchaseHandler(c *fiber.Ctx) error
	CreateDeploymentHandler(c *fiber.Ctx) error
	GetInventoryStatusHandler(c *fiber.Ctx) error
}

type InventoryHandler struct {
	service InventoryServiceInterface
}

func NewInventoryHandler(service InventoryServiceInterface) InventoryHandlerInterface {
	return &InventoryHandler{service: service}
}

// CreatePurchaseHandler handles POST /api/purchases
func (h *InventoryHandler) CreatePurchaseHandler(c *fiber.Ctx) error {
	request := CreatePurchaseRequest{}

	if err := c.BodyParser(&request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Invalid request body", nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "User ID not found in token", nil)
	}

	response, err := h.service.CreatePurchaseService(request, userID)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Purchase created successfully", response)
}

// CreateDeploymentHandler handles POST /api/deployments
func (h *InventoryHandler) CreateDeploymentHandler(c *fiber.Ctx) error {
	request := CreateDeploymentRequest{}

	if err := c.BodyParser(&request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, "Invalid request body", nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return helpers.ResponseUtils(c, fiber.StatusUnauthorized, false, "User ID not found in token", nil)
	}

	response, err := h.service.CreateDeploymentService(request, userID)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Asset deployment successful", response)
}

// GetInventoryStatusHandler handles GET /api/inventory/status
func (h *InventoryHandler) GetInventoryStatusHandler(c *fiber.Ctx) error {
	request := InventoryStatusRequest{}

	// Parse query parameters
	if brand := c.Query("brand"); brand != "" {
		request.Brand = &brand
	}
	if model := c.Query("model"); model != "" {
		request.Model = &model
	}
	if status := c.Query("status"); status != "" {
		request.Status = &status
	}

	response, err := h.service.GetInventoryStatusService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Inventory status retrieved successfully", response)
}
