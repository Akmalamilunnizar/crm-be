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
		AssignedTo:  &uid,
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

	// Parse multipart form (note + optional image)
	log.Printf("SendToNOC Handler: Parsing multipart form for ticket %d", id)
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("SendToNOC Handler: Failed to parse multipart form: %v", err)
		return helpers.ResponseUtils(c, 400, false, "failed to parse multipart form", nil)
	}

	log.Printf("SendToNOC Handler: Form values: %v", form.Value)
	log.Printf("SendToNOC Handler: Form files: %v", form.File)

	notes := form.Value["note"]
	if len(notes) == 0 {
		log.Printf("SendToNOC Handler: No note provided in form")
		return helpers.ResponseUtils(c, 400, false, "note is required", nil)
	}
	note := notes[0]
	log.Printf("SendToNOC Handler: Note received: %s", note)

	var imageFilename *string
	if files := form.File["image"]; len(files) > 0 {
		log.Printf("SendToNOC Handler: Processing image file")
		file := files[0]
		log.Printf("SendToNOC Handler: Image file name: %s, size: %d", file.Filename, file.Size)
		if file.Size > 10*1024*1024 {
			return helpers.ResponseUtils(c, 400, false, "image file too large (max 10MB)", nil)
		}
		uploadDir := "uploads/cs-images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-cs%s", id, ext)
		fullPath := filepath.Join(uploadDir, filename)
		log.Printf("SendToNOC Handler: Saving image to: %s", fullPath)
		if err := c.SaveFile(file, fullPath); err != nil {
			log.Printf("SendToNOC Handler: Failed to save image: %v", err)
			return helpers.ResponseUtils(c, 500, false, "failed to save image file", nil)
		}
		log.Printf("SendToNOC Handler: Image saved successfully")
		imageFilename = &filename
	} else {
		log.Printf("SendToNOC Handler: No image file provided")
	}

	out, err := h.svc.SendToNOC(id, note, imageFilename)
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

	// Get type from form
	var tType *string
	if types := form.Value["type"]; len(types) > 0 && types[0] != "" {
		tType = &types[0]
	}

	// Get uploaded image file
	var imageFilename *string

	// Handle image file
	if files := form.File["image"]; len(files) > 0 {
		file := files[0]
		if file.Size > 10*1024*1024 { // 10MB limit
			return helpers.ResponseUtils(c, 400, false, "image file too large (max 10MB)", nil)
		}

		// Create uploads directory if not exists
		uploadDir := "uploads/noc-images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}

		// Generate filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-noc%s", id, ext)
		filepath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveFile(file, filepath); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to save image file", nil)
		}

		// Store filename for database
		imageFilename = &filename
	}

	out, err := h.svc.SendToCS(id, note, tType, imageFilename)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "Ticket sent to CS successfully", fiber.Map{
		"ticket_id":        out.ID,
		"note":             out.NOCNote,
		"type":             out.Type,
		"img_noc":          out.ImgNOC,
		"current_assignee": out.CurrentAssignee,
	})
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
	// No longer require technician_id - assign to TECHNICIAN role in general
	out, err := h.svc.AssignTechnician(id)
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

func (h *Handler) CSResolve(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}
	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	out, err := h.svc.CSResolve(id, body.Note)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ok", out)
}

// AddTechnicianNote allows technician to add a note to a ticket
func (h *Handler) AddTechnicianNote(c *fiber.Ctx) error {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("AddTechnicianNote: PANIC recovered: %v", r)
		}
	}()

	log.Printf("AddTechnicianNote: Starting request for ticket ID: %s", c.Params("id"))

	id, err := idParam(c)
	if err != nil {
		log.Printf("AddTechnicianNote: Invalid ID parameter: %v", err)
		return helpers.ResponseUtils(c, 400, false, "bad id", nil)
	}

	log.Printf("AddTechnicianNote: Parsing multipart form for ticket ID: %d", id)
	log.Printf("AddTechnicianNote: Content-Type header: %s", c.Get("Content-Type"))

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("AddTechnicianNote: Failed to parse multipart form: %v", err)
		log.Printf("AddTechnicianNote: Request body size: %d", len(c.Body()))
		return helpers.ResponseUtils(c, 400, false, "failed to parse multipart form", nil)
	}

	log.Printf("AddTechnicianNote: Successfully parsed multipart form")

	// Get note from form
	notes := form.Value["note"]
	if len(notes) == 0 || notes[0] == "" {
		log.Printf("AddTechnicianNote: Note is required but not provided")
		return helpers.ResponseUtils(c, 400, false, "note is required", nil)
	}
	note := notes[0]
	log.Printf("AddTechnicianNote: Note received: %s", note)

	// Get uploaded files
	var imgTechBfPath, imgTechAfPath *string

	log.Printf("AddTechnicianNote: Checking for uploaded files")
	log.Printf("AddTechnicianNote: Available form fields: %v", form.Value)
	log.Printf("AddTechnicianNote: Available file fields: %v", form.File)

	// Handle img_tech_bf file - using exact same logic as SendToCS
	if files := form.File["img_tech_bf"]; len(files) > 0 {
		log.Printf("AddTechnicianNote: Processing img_tech_bf file")
		file := files[0]
		log.Printf("AddTechnicianNote: img_tech_bf file size: %d bytes", file.Size)
		log.Printf("AddTechnicianNote: img_tech_bf filename: %s", file.Filename)

		if file.Size > 10*1024*1024 { // 10MB limit
			log.Printf("AddTechnicianNote: img_tech_bf file too large: %d bytes", file.Size)
			return helpers.ResponseUtils(c, 400, false, "img_tech_bf file too large (max 10MB)", nil)
		}

		// Create uploads directory if not exists - using separate directory
		uploadDir := "uploads/technician-images"
		log.Printf("AddTechnicianNote: Creating upload directory: %s", uploadDir)
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("AddTechnicianNote: Failed to create upload directory: %v", err)
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}

		// Generate filename - using same pattern as SendToCS
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-bf%s", id, ext)
		filepath := filepath.Join(uploadDir, filename)
		log.Printf("AddTechnicianNote: Saving img_tech_bf to: %s", filepath)

		// Save file - using exact same logic as SendToCS
		if err := c.SaveFile(file, filepath); err != nil {
			log.Printf("AddTechnicianNote: Failed to save img_tech_bf file: %v", err)
			return helpers.ResponseUtils(c, 500, false, "failed to save img_bf file", nil)
		}

		log.Printf("AddTechnicianNote: Successfully saved img_tech_bf file")

		// Store filename for database
		imgTechBfPath = &filename
	}

	// Handle img_tech_af file - using exact same logic as SendToCS
	if files := form.File["img_tech_af"]; len(files) > 0 {
		file := files[0]
		if file.Size > 10*1024*1024 { // 10MB limit
			return helpers.ResponseUtils(c, 400, false, "img_tech_af file too large (max 10MB)", nil)
		}

		// Create uploads directory if not exists - using separate directory
		uploadDir := "uploads/technician-images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
		}

		// Generate filename - using same pattern as SendToCS
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d-af%s", id, ext)
		filepath := filepath.Join(uploadDir, filename)

		// Save file - using exact same logic as SendToCS
		if err := c.SaveFile(file, filepath); err != nil {
			return helpers.ResponseUtils(c, 500, false, "failed to save img_tech_af file", nil)
		}

		// Store filename for database
		imgTechAfPath = &filename
	}

	log.Printf("AddTechnicianNote: Calling service with note: %s, imgTechBfPath: %v, imgTechAfPath: %v", note, imgTechBfPath, imgTechAfPath)

	out, err := h.svc.AddTechnicianNote(id, note, imgTechBfPath, imgTechAfPath)
	if err != nil {
		log.Printf("AddTechnicianNote: Service error: %v", err)
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	log.Printf("AddTechnicianNote: Service call successful, returning response")
	return helpers.ResponseUtils(c, 200, true, "Technician note & images added successfully", fiber.Map{
		"ticket_id":       out.ID,
		"technician_note": out.TechnicianNote,
		"img_tech_bf":     out.ImgTechBF,
		"img_tech_af":     out.ImgTechAF,
		"img_cs":          out.ImgCS,
		"img_noc":         out.ImgNOC,
	})
}

func (h *Handler) ReportByType(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	rows, err := h.svc.repo.CountByType(startDate, endDate)
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

func (h *Handler) UploadNOCImage(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("image")
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "no image file provided", nil)
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return helpers.ResponseUtils(c, 400, false, "invalid file type. only images are allowed", nil)
	}

	// Validate file size (10MB limit)
	if file.Size > 10*1024*1024 {
		return helpers.ResponseUtils(c, 400, false, "file size too large. maximum 10MB allowed", nil)
	}

	// Generate unique filename
	filename := generateUniqueFilename(file.Filename)

	// Save file to uploads directory
	uploadDir := "./uploads/noc-images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return helpers.ResponseUtils(c, 500, false, "failed to create upload directory", nil)
	}

	filepath := uploadDir + "/" + filename
	if err := c.SaveFile(file, filepath); err != nil {
		return helpers.ResponseUtils(c, 500, false, "failed to save file", nil)
	}

	return helpers.ResponseUtils(c, 200, true, "image uploaded successfully", map[string]string{
		"filename": filename,
	})
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	randomStr := fmt.Sprintf("%d", timestamp)
	return randomStr + ext
}
