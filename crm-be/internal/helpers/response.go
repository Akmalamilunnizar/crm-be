package helpers

import "github.com/gofiber/fiber/v2"

func ResponseUtils(c *fiber.Ctx, statusCode int, success bool, message string, data any) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": success,
		"message": message,
		"data":    data,
	})
}
