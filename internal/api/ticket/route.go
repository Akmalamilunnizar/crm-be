package ticketapi

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func TicketRoutes(api fiber.Router) {
	log.Println("TicketRoutes: Starting route registration...")

	db := database.GetDB()
	r := NewRepo(db)
	s := NewService(r)
	h := NewHandler(s)

	g := api.Group("/tickets")
	log.Println("TicketRoutes: Created /tickets group")

	g.Get("/", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.List)
	g.Post("/", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.Create)

	g.Post("/:id/send-to-noc", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.SendToNOC)
	g.Post("/:id/send-to-cs", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "NOC", "TECHNICIAN"), h.SendToCS)
	g.Post("/:id/noc-solved", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "NOC"), h.NOCSolved)
	g.Post("/:id/noc-physical", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "NOC"), h.NOCPhysical)
	g.Post("/:id/assign-technician", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.AssignTechnician)
	g.Post("/:id/resolve", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "CUSTOMER SERVICE"), h.CSResolve)
	g.Post("/:id/technician-resolve", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.TechnicianResolve)
	g.Post("/:id/technician-note", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.AddTechnicianNote)
	g.Post("/upload-noc-image", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "NOC"), h.UploadNOCImage)

	g.Get("/reports/by-type", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.ReportByType)
	g.Get("/lookups/trouble-types", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.TroubleTypes)
	g.Post("/lookups/trouble-types", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.CreateTroubleType)
	g.Get("/reports/hotspots", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.HotLocations)
	// Polling endpoint for updates
	g.Get("/updates", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.UpdatesSince)

	// Debug endpoint to test role extraction
	g.Get("/debug/role", helpers.VerifyToken, func(c *fiber.Ctx) error {
		log.Println("Debug endpoint accessed")
		role := c.Locals("role").(string)
		userID := c.Locals("user_id").(string)

		// Use the same normalization as RequireRoles
		normalizedRole := helpers.NormalizeRole(role)

		return helpers.ResponseUtils(c, 200, true, "debug", fiber.Map{
			"role":            role,
			"user_id":         userID,
			"normalized_role": normalizedRole,
		})
	})

	log.Println("TicketRoutes: All routes registered successfully")
}
