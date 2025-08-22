package transaction

import (
	"skripsi-be/internal/api/common/validation"
	"skripsi-be/internal/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AdminTransactionHandlerStruct struct {
	service AdminTransactionServiceInterface
}

func NewAdminTransactionHandler(service AdminTransactionServiceInterface) AdminTransactionHandlerStruct {
	return AdminTransactionHandlerStruct{service}
}

func (h AdminTransactionHandlerStruct) GetAllAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	request := SearchAdminTransactionRequest{}
	request.TypeCash = c.Query("type_cash")
	request.Type = c.Query("type")
	transactions, err := h.service.GetAllAdminTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", transactions)
}

func (h AdminTransactionHandlerStruct) GetByIdAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	request := IdAdminTransactionRequest{}
	request.Id = c.Params("id")
	transaction, err := h.service.GetByIdAdminTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", transaction)

}

func (h AdminTransactionHandlerStruct) CreateAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	request := CreateAdminTransactionRequest{}
	c.BodyParser(&request)

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}
	transaction, err := h.service.CreateAdminTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", transaction)
}

func (h AdminTransactionHandlerStruct) UpdateAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	request := UpdateAdminTransactionRequest{}

	c.BodyParser(&request)
	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}
	transaction, err := h.service.UpdateAdminTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", transaction)
}

func (h AdminTransactionHandlerStruct) DeleteAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	request := IdAdminTransactionRequest{}

	request.Id = c.Params("id")

	errValidation := validation.ValidationRequest(&request)
	if errValidation != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, strings.Join(errValidation, ", "), "")
	}

	transaction, err := h.service.DeleteAdminTransactionService(request)
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}

	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", transaction)
}

func (h AdminTransactionHandlerStruct) SyncAdminTransactionHandlerStruct(c *fiber.Ctx) error {
	err := h.service.SyncAdminTransactionService()

	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), "")
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", "Sync Success")
}
