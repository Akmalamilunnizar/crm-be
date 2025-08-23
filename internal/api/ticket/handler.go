package ticketapi

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// AddTechnicianNote allows technician to add a note to a ticket
func (h *Handler) AddTechnicianNote(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "failed to parse multipart form", nil)
	}

	// Get note from form
	notes := form.Value["note"]
	if len(notes) == 0 || notes[0] == "" {
		return helpers.ResponseUtils(c, 400, false, "note is required", nil)
	}
	note := notes[0]

	// Get uploaded files
	var imgTechBfPath, imgTechAfPath *string

	// Handle img_tech_bf file
	if files := form.File["img_tech_bf"]; len(files) > 0 {
		file := files[0]
		if file.Size > 10*1024*1024 { // 10MB limit
			return helpers.ResponseUtils(c, 400, false, "img_tech_bf file too large (max 10MB)", nil)
		}

		// Create uploads directory if not exists
		uploadDir := "uploads/technician"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}

		// Generate filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-bf%s", id, ext)
		filepath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveFile(file, filepath); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to save img_bf file", nil)
		}

		// Store relative path for database
		relativePath := "/" + filepath
		imgTechBfPath = &relativePath
	}

	// Handle img_tech_af file
	if files := form.File["img_tech_af"]; len(files) > 0 {
		file := files[0]
		if file.Size > 10*1024*1024 { // 10MB limit
			return helpers.ResponseUtils(c, 400, false, "img_tech_af file too large (max 10MB)", nil)
		}

		// Create uploads directory if not exists
		uploadDir := "uploads/technician"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}

		// Generate filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-af%s", id, ext)
		filepath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveFile(file, filepath); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to save img_tech_af file", nil)
		}

		// Store relative path for database
		relativePath := "/" + filepath
		imgTechAfPath = &relativePath
	}

	out, err := h.svc.AddTechnicianNote(id, note, imgTechBfPath, imgTechAfPath)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Technician note & images added successfully", fiber.Map{
		"ticket_id":       out.ID,
		"technician_note": out.TechnicianNote,
		"img_tech_bf":     out.ImgTechBf,
		"img_tech_af":     out.ImgTechAf,
		"img_cs":          out.ImgCs,
		"img_noc":         out.ImgNoc,
	})
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
