package report

import (
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type AdminReportHandlerStruct struct {
	service AdminReportServiceInterface
}

func NewAdminReportHandler(service AdminReportServiceInterface) *AdminReportHandlerStruct {
	return &AdminReportHandlerStruct{service}
}

func (h AdminReportHandlerStruct) GetAllAdminReportHandler(c *fiber.Ctx) error {
	customers, err := h.service.GetAllInternetAdminReportService()
	if err != nil {
		return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, fiber.StatusOK, true, "", customers)
}
