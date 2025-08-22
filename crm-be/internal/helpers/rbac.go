package helpers

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Role mapping to handle different role name formats
var roleMapping = map[string]string{
	"CUSTOMER SERVICE": "CUSTOMER_SERVICE",
	"CUSTOMER_SERVICE": "CUSTOMER_SERVICE",
	"ADMIN":            "ADMIN",
	"NOC":              "NOC",
	"TECHNICIAN":       "TECHNICIAN",
	"FINANCE":          "FINANCE",
}

func NormalizeRole(role string) string {
	if normalized, exists := roleMapping[role]; exists {
		return normalized
	}
	// Fallback: replace spaces with underscores
	return strings.ReplaceAll(role, " ", "_")
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		v := c.Locals("role")
		role, _ := v.(string)

		// Normalize role name using comprehensive mapping
		normalizedRole := NormalizeRole(role)

		// Debug logging
		log.Printf("RequireRoles - Original role: '%s', Normalized role: '%s'", role, normalizedRole)
		log.Printf("RequireRoles - Allowed roles: %v", roles)

		if _, ok := allowed[normalizedRole]; !ok {
			log.Printf("RequireRoles - Access denied for role: '%s' (normalized: '%s')", role, normalizedRole)
			return ResponseUtils(c, fiber.StatusForbidden, false, "forbidden", nil)
		}

		log.Printf("RequireRoles - Access granted for role: '%s' (normalized: '%s')", role, normalizedRole)
		return c.Next()
	}
}
