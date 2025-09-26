package ticketapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"skripsi-be/internal/helpers"
	"skripsi-be/internal/models/entities"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTechnicianSteps gets all predefined technician steps
func (h *Handler) GetTechnicianSteps(c *fiber.Ctx) error {
	steps, err := h.svc.GetTechnicianSteps()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "technician steps retrieved", steps)
}

// GetSpareParts gets all available spare parts
func (h *Handler) GetSpareParts(c *fiber.Ctx) error {
	parts, err := h.svc.GetSpareParts()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "spare parts retrieved", parts)
}

// GetTicketTechnicianSteps gets technician progress for a specific ticket
func (h *Handler) GetTicketTechnicianSteps(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	technicianID := c.Query("technician_id")
	if technicianID == "" || technicianID == "TECHNICIAN" || technicianID == "ADMIN" || technicianID == "CUSTOMER_SERVICE" {
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			technicianID = uid
		}
	}

	steps, err := h.svc.GetTicketTechnicianSteps(id, technicianID)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "ticket technician steps retrieved", steps)
}

// UpdateTechnicianStep updates a technician's progress on a specific step
func (h *Handler) UpdateTechnicianStep(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	// Handle multipart form data for images
	form, err := c.MultipartForm()
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "failed to parse multipart form", nil)
	}

	var req struct {
		StepID         uint64  `form:"step_id"`
		TechnicianID   string  `form:"technician_id"`
		Status         string  `form:"status"`
		Notes          *string `form:"notes"`
		SparePartsUsed *string `form:"spare_parts_used"`
	}

	// Parse form fields
	if stepIDStr := form.Value["step_id"]; len(stepIDStr) > 0 {
		if stepID, err := strconv.ParseUint(stepIDStr[0], 10, 64); err == nil {
			req.StepID = stepID
		}
	}
	if technicianID := form.Value["technician_id"]; len(technicianID) > 0 {
		req.TechnicianID = technicianID[0]
	}
	if status := form.Value["status"]; len(status) > 0 {
		req.Status = status[0]
	}
	if notes := form.Value["notes"]; len(notes) > 0 && notes[0] != "" {
		req.Notes = &notes[0]
	}
	if spareParts := form.Value["spare_parts_used"]; len(spareParts) > 0 && spareParts[0] != "" {
		req.SparePartsUsed = &spareParts[0]
	}

	// Resolve technician ID from JWT if missing or looks like a role label
	resolvedTechID := req.TechnicianID
	if resolvedTechID == "" || resolvedTechID == "TECHNICIAN" || resolvedTechID == "ADMIN" || resolvedTechID == "CUSTOMER_SERVICE" {
		if auth := c.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			parts := strings.Split(strings.TrimPrefix(auth, "Bearer "), ".")
			if len(parts) == 3 {
				if payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1]); err == nil {
					var payload map[string]interface{}
					if err := json.Unmarshal(payloadBytes, &payload); err == nil {
						if sub, ok := payload["sub"].(string); ok && sub != "" {
							resolvedTechID = sub
						}
					}
				}
			}
		}
	}

	// Handle uploaded images
	var imagePaths []string
	// Ensure upload directory exists (under public so it can be served statically)
	uploadDir := filepath.Join("public", "uploads", "technician-progress")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		_ = os.MkdirAll(uploadDir, 0755)
	}
	if images := form.File["images"]; len(images) > 0 {
		for _, file := range images {
			// Save file to public/uploads/technician-progress/
			filename := fmt.Sprintf("%d_%s_%s", time.Now().Unix(), resolvedTechID, file.Filename)
			fullpath := filepath.Join(uploadDir, filename)

			if err := c.SaveFile(file, fullpath); err != nil {
				return helpers.ResponseUtils(c, 500, false, fmt.Sprintf("failed to save image: %v", err), nil)
			}

			imagePaths = append(imagePaths, filename)
		}
	}

	// Persist JSON image paths in step record
	var imagePathsJSON *string
	if len(imagePaths) > 0 {
		b, _ := json.Marshal(imagePaths)
		s := string(b)
		imagePathsJSON = &s
	}

	// Basic required fields validation
	if req.StepID == 0 || resolvedTechID == "" || req.Status == "" {
		return helpers.ResponseUtils(c, 400, false, "step_id, technician_id, and status are required", nil)
	}

	if err := h.svc.UpdateTechnicianStep(id, req.StepID, resolvedTechID, req.Status, req.Notes, req.SparePartsUsed, imagePathsJSON); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	return helpers.ResponseUtils(c, 200, true, "technician step updated successfully", nil)
}

// SetNetworkArchitecture sets FTTH/HTB for a ticket
func (h *Handler) SetNetworkArchitecture(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	var req struct {
		Architecture string `json:"architecture"`
	}
	if err := c.BodyParser(&req); err != nil || req.Architecture == "" {
		return helpers.ResponseUtils(c, 400, false, "architecture is required", nil)
	}
	if err := h.svc.SetNetworkArchitecture(id, req.Architecture); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "network architecture updated", nil)
}

// UpsertTechnicianTeam assigns Senior/Junior/Helper ensuring unique users
func (h *Handler) UpsertTechnicianTeam(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	var req struct {
		Members []struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		} `json:"members"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseUtils(c, 400, false, "invalid body", nil)
	}
	members := make([]entities.TechnicianTeamMember, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, entities.TechnicianTeamMember{TicketID: id, UserID: m.UserID, Role: strings.ToLower(m.Role)})
	}
	if err := h.svc.UpsertTechnicianTeam(id, members); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "technician team updated", nil)
}

// GetTechnicianStepProgress gets progress summary for a technician on a ticket
func (h *Handler) GetTechnicianStepProgress(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	technicianID := c.Query("technician_id")
	if technicianID == "" || technicianID == "TECHNICIAN" || technicianID == "ADMIN" || technicianID == "CUSTOMER_SERVICE" {
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			technicianID = uid
		}
	}

	progress, err := h.svc.GetTechnicianStepProgress(id, technicianID)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	return helpers.ResponseUtils(c, 200, true, "technician step progress retrieved", progress)
}

// MarkTechnicianJobCompleted marks the entire technician job as completed
func (h *Handler) MarkTechnicianJobCompleted(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	if err := h.svc.MarkTechnicianJobCompleted(id); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	// Return updated ticket so clients can refresh UI status
	ticket, _ := h.svc.repo.GetTicketByID(id)
	return helpers.ResponseUtils(c, 200, true, "technician job marked as completed", ticket)
}

// GetTechnicianChecklist gets the complete checklist for a technician on a ticket
func (h *Handler) GetTechnicianChecklist(c *fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}

	technicianID := c.Query("technician_id")
	if technicianID == "" || technicianID == "TECHNICIAN" || technicianID == "ADMIN" || technicianID == "CUSTOMER_SERVICE" {
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			technicianID = uid
		}
	}

	// Get all predefined steps
	steps, err := h.svc.GetTechnicianSteps()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	// Get technician progress
	progress, err := h.svc.GetTicketTechnicianSteps(id, technicianID)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}

	// Create a map of progress for quick lookup
	progressMap := make(map[uint64]entities.TicketTechnicianStep)
	for _, p := range progress {
		progressMap[p.StepID] = p
	}

	// Combine steps with progress
	var checklist []map[string]interface{}
	for _, step := range steps {
		item := map[string]interface{}{
			"step_id":          step.ID,
			"step_order":       step.StepOrder,
			"title":            step.Title,
			"description":      step.Description,
			"tools":            step.Tools,
			"spare_parts":      step.SpareParts,
			"procedure":        step.Procedure,
			"solution":         step.Solution,
			"status":           "pending",
			"notes":            nil,
			"spare_parts_used": nil,
			"completed_at":     nil,
			"image_paths":      nil,
		}

		// Add progress if exists
		if p, exists := progressMap[step.ID]; exists {
			item["status"] = p.Status
			item["notes"] = p.Notes
			item["spare_parts_used"] = p.SparePartsUsed
			item["completed_at"] = p.CompletedAt
			item["image_paths"] = p.ImagePaths
		}

		checklist = append(checklist, item)
	}

	// also include ticket network architecture to drive UI
	ticket, terr := h.svc.repo.GetTicketByID(id)
	if terr != nil {
		// do not fail the whole response; return checklist only
		return helpers.ResponseUtils(c, 200, true, "technician checklist retrieved", fiber.Map{
			"checklist":            checklist,
			"network_architecture": nil,
		})
	}

	return helpers.ResponseUtils(c, 200, true, "technician checklist retrieved", fiber.Map{
		"checklist":            checklist,
		"network_architecture": ticket.NetworkArchitecture,
	})
}
