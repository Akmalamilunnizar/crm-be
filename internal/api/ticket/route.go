package ticketapi

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/services"

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
	g.Post("/lookups/trouble-types", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC"), h.CreateTroubleType)

	// ML Classification endpoints
	g.Post("/classify", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.ClassifyTicket)
	g.Get("/ml/stats", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.GetMLStats)
	g.Get("/reports/hotspots", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.HotLocations)
	// Polling endpoint for updates
	g.Get("/updates", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), h.UpdatesSince)

	// Technician workflow
	g.Post("/:id/accept", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.Accept)
	// Consolidated team assignment endpoint with validation
	g.Post("/:id/team", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.UpsertTechnicianTeam)
	// Fetch current team members for a ticket
	g.Get("/:id/team-members", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN", "CUSTOMER_SERVICE"), h.GetTechnicianTeamMembers)
	g.Post("/:id/steps", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.AddStep)
	g.Post("/:id/technician-completed", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.MarkTechnicianCompleted)

	// Network architecture
	g.Post("/:id/network-architecture", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.SetNetworkArchitecture)

	// Technician workflow routes
	g.Get("/technician-steps", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.GetTechnicianSteps)
	g.Get("/spare-parts", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.GetSpareParts)
	g.Get("/:id/technician-checklist", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN", "CUSTOMER_SERVICE"), h.GetTechnicianChecklist)
	g.Get("/:id/technician-steps", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN", "CUSTOMER_SERVICE"), h.GetTicketTechnicianSteps)
	g.Post("/:id/technician-step", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.UpdateTechnicianStep)
	g.Post("/:id/technician-selfie", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.SaveSelfieStep)
	g.Get("/:id/technician-selfie", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN", "CUSTOMER_SERVICE"), h.GetSelfieStep)
	g.Get("/:id/technician-progress", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN", "CUSTOMER_SERVICE"), h.GetTechnicianStepProgress)
	g.Post("/:id/technician-complete", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "TECHNICIAN"), h.MarkTechnicianJobCompleted)

	// CS verification/close
	g.Post("/:id/verify-close", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), h.VerifyAndClose)

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

	// Accumulation management routes
	log.Println("TicketRoutes: Registering accumulation routes...")
	log.Println("TicketRoutes: Creating accumulation service...")
	accumulationService := services.NewAccumulationService(db)
	log.Println("TicketRoutes: Accumulation service created successfully")
	log.Println("TicketRoutes: Creating accumulation handler...")
	accumulationHandler := NewAccumulationHandler(accumulationService)
	log.Println("TicketRoutes: Accumulation handler created successfully")

	// Accumulation endpoints
	g.Get("/:id/similar", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), accumulationHandler.GetSimilarTroubles)
	g.Post("/accumulation", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), accumulationHandler.UpdateAccumulation)
	g.Post("/accumulation/auto-detect", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE"), accumulationHandler.AutoDetectAndGroup)
	g.Get("/accumulation/stats", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), accumulationHandler.GetAccumulationStats)
	g.Get("/accumulation/high", helpers.VerifyToken, helpers.RequireRoles("ADMIN", "CUSTOMER_SERVICE", "NOC", "TECHNICIAN"), accumulationHandler.GetHighAccumulationTickets)

	log.Println("TicketRoutes: Accumulation routes registered successfully")

	log.Println("TicketRoutes: All routes registered successfully")
}
