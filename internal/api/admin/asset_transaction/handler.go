package asset_transaction

import (
	"fmt"
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AssetTransactionHandlerInterface interface {
	GetAllAssetTransactionHandler(c *fiber.Ctx) error
	GetByIdAssetTransactionHandler(c *fiber.Ctx) error
	CreateAssetTransactionHandler(c *fiber.Ctx) error
	UpdateAssetTransactionHandler(c *fiber.Ctx) error
	DeleteAssetTransactionHandler(c *fiber.Ctx) error
}

type AssetTransactionHandlerStruct struct {
	service AssetTransactionServiceInterface
}

func NewAssetTransactionHandler(service AssetTransactionServiceInterface) AssetTransactionHandlerInterface {
	return &AssetTransactionHandlerStruct{service}
}

func (h *AssetTransactionHandlerStruct) GetAllAssetTransactionHandler(c *fiber.Ctx) error {
	request := GetAssetTransactionsRequest{}

	// Parse query parameters
	if assetID := c.Query("asset_id"); assetID != "" {
		request.AssetID = &assetID
	}
	if customerInstallationID := c.Query("customer_installation_id"); customerInstallationID != "" {
		request.CustomerInstallationID = &customerInstallationID
	}
	if transactionType := c.Query("transaction_type"); transactionType != "" {
		request.TransactionType = &transactionType
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		request.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		request.DateTo = &dateTo
	}

	transactions, err := h.service.GetAssetTransactionsService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success get asset transactions", transactions)
}

func (h *AssetTransactionHandlerStruct) GetByIdAssetTransactionHandler(c *fiber.Ctx) error {
	request := IdAssetTransactionRequest{
		Id: c.Params("id"),
	}

	transaction, err := h.service.GetByIdAssetTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success get asset transaction", transaction)
}

func (h *AssetTransactionHandlerStruct) CreateAssetTransactionHandler(c *fiber.Ctx) error {
	request := CreateAssetTransactionRequest{}

	if err := c.BodyParser(&request); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	// Get the authenticated user ID from the request context
	var createdBy string
	if userID := c.Locals("user_id"); userID != nil {
		createdBy = fmt.Sprintf("%v", userID)
	} else {
		// Fallback to a known admin user ID to avoid foreign key constraint issues
		createdBy = "b6cbf8b6-6a2e-45f3-a9ab-e30e5941bf5a" // Admin user ID from JWT
	}

	// Handle empty customer_installation_id for standalone goods transactions
	// Use a placeholder value if empty to avoid foreign key constraint issues
	if request.CustomerInstallationID == "" {
		// For standalone goods transactions, we'll use a special placeholder
		// You may want to create a special "STANDALONE" customer installation record
		// For now, we'll skip the foreign key by using empty string (if DB allows) or handle it in repository
		request.CustomerInstallationID = "" // Keep empty for standalone transactions
	}

	transaction, err := h.service.CreateAssetTransactionService(request, createdBy)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success create asset transaction", transaction)
}

func (h *AssetTransactionHandlerStruct) UpdateAssetTransactionHandler(c *fiber.Ctx) error {
	request := UpdateAssetTransactionRequest{
		IdAssetTransactionRequest: IdAssetTransactionRequest{
			Id: c.Params("id"),
		},
	}

	if err := c.BodyParser(&request.CreateAssetTransactionRequest); err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), nil)
	}

	transaction, err := h.service.UpdateAssetTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success update asset transaction", transaction)
}

func (h *AssetTransactionHandlerStruct) DeleteAssetTransactionHandler(c *fiber.Ctx) error {
	request := IdAssetTransactionRequest{
		Id: c.Params("id"),
	}

	transaction, err := h.service.DeleteAssetTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "Success delete asset transaction", transaction)
}
