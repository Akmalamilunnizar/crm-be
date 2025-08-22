package ticketapi

import (
	"log"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc *Service }

func NewHandler(s *Service) *Handler { return &Handler{s} }

func (h *Handler) List(c *fiber.Ctx) error {
	log.Printf("Ticket List handler called")

	// Debug: Check roles and assignments
	h.svc.repo.DebugRolesAndAssignments()

	items, err := h.svc.repo.ListAll()
	if err != nil {
		log.Printf("Ticket List error: %v", err)
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	log.Printf("Ticket List returned %d items", len(items))
	return helpers.ResponseUtils(c, 200, true, "ok", items)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var in struct {
		CustomerID  string  `json:"customer_id"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Type        *string `json:"type"`
	}
	if err := c.BodyParser(&in); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	// assigned_to must reference users.id → use authenticated user id from token
	uidVal := c.Locals("user_id")
	uid, _ := uidVal.(string)
	if uid == "" {
		return helpers.ResponseUtils(c, 401, false, "unauthorized", nil)
	}
	t := entities.TroubleTicket{
		CustomerID:  in.CustomerID,
		Title:       in.Title,
		Description: in.Description,
		Type:        in.Type,
		AssignedTo:  uid,
	}
	out, err := h.svc.CreateCS(t)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 201, true, "created", out)
}

func idParam(c *fiber.Ctx) (uint64, error) { return strconv.ParseUint(c.Params("id"), 10, 64) }

func (h *Handler) SendToNOC(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	out, err := h.svc.SendToNOC(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

func (h *Handler) NOCSolved(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	out, err := h.svc.NOCSolved(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

// SendToCS allows NOC (or ADMIN) to send the ticket back to Customer Service with a note
func (h *Handler) SendToCS(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	out, err := h.svc.SendToCS(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

func (h *Handler) NOCPhysical(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	out, err := h.svc.NOCPhysical(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

func (h *Handler) AssignTechnician(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		TechnicianID string `json:"technician_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	out, err := h.svc.AssignTechnician(id, body.TechnicianID)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

func (h *Handler) TechnicianResolve(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	out, err := h.svc.TechnicianResolve(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

func (h *Handler) ReportByType(c *fiber.Ctx) error {
	rows, err := h.svc.repo.CountByType()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", rows)
}

// GET trouble types lookup
func (h *Handler) TroubleTypes(c *fiber.Ctx) error {
	log.Printf("TroubleTypes handler called")

	rows, err := h.svc.repo.AllTroubleTypes()
	if err != nil {
		log.Printf("TroubleTypes error: %v", err)
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	log.Printf("TroubleTypes returned %d items", len(rows))
	return helpers.ResponseUtils(c, 200, true, "ok", rows)
}

// GET hotspots by customer geolocation
func (h *Handler) HotLocations(c *fiber.Ctx) error {
	rows, err := h.svc.repo.HotLocations()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", rows)
}

// POST create trouble type
func (h *Handler) CreateTroubleType(c *fiber.Ctx) error {
	var in struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := c.BodyParser(&in); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	troubleType := entities.TroubleTypeRow{
		ID:   in.ID,
		Name: &in.Name,
	}

	err := h.svc.repo.CreateTroubleType(troubleType)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 201, true, "created", troubleType)
}

// GET /api/tickets/updates?since=RFC3339
func (h *Handler) UpdatesSince(c *fiber.Ctx) error {
	sinceStr := c.Query("since")
	if sinceStr == "" {
		// default to 60 seconds ago
		since := time.Now().Add(-60 * time.Second)
		sinceStr = since.Format(time.RFC3339)
	}
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad since timestamp", nil)
	}

	roleVal := c.Locals("role")
	role, _ := roleVal.(string)
	normalized := helpers.NormalizeRole(role)

	uidVal := c.Locals("user_id")
	userID, _ := uidVal.(string)

	items, err := h.svc.repo.UpdatesSinceForRole(since, normalized, userID)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", items)
}
