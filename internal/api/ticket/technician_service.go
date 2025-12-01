package ticketapi

import (
	"encoding/json"
	"fmt"
	"skripsi-be/internal/models/entities"
)

// GetTechnicianSteps gets all predefined technician steps
// If networkArchitecture is provided, filters steps to only those applicable to that architecture
func (s *Service) GetTechnicianSteps(networkArchitecture *string) ([]entities.TechnicianStep, error) {
	return s.repo.GetTechnicianSteps(networkArchitecture)
}

// GetSpareParts gets all available spare parts
func (s *Service) GetSpareParts() ([]entities.SparePart, error) {
	return s.repo.GetSpareParts()
}

// GetTicketTechnicianSteps gets technician progress for a specific ticket
func (s *Service) GetTicketTechnicianSteps(ticketID uint64, technicianID string) ([]entities.TicketTechnicianStep, error) {
	return s.repo.GetTicketTechnicianSteps(ticketID, technicianID)
}

// UpdateTechnicianStep updates a technician's progress on a specific step
func (s *Service) UpdateTechnicianStep(ticketID uint64, stepID uint64, technicianID string, status string, notes *string, sparePartsUsed *string, imagePaths *string) error {
	// Validate status
	validStatuses := []string{"pending", "done", "not_applicable", "needs_spare_parts"}
	valid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid status: %s. Must be one of: %v", status, validStatuses)
	}
	// Architecture must be selected first
	t, err := s.repo.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if t.NetworkArchitecture == nil || *t.NetworkArchitecture == "" {
		return fmt.Errorf("network architecture must be selected before progressing")
	}

	// Enforce sequential steps
	st, err := s.repo.GetStepByID(stepID)
	if err != nil {
		return err
	}
	hasIncomplete, err := s.repo.HasIncompletePreviousSteps(ticketID, technicianID, st.StepOrder)
	if err != nil {
		return err
	}
	if hasIncomplete {
		return fmt.Errorf("previous steps must be completed before this step")
	}

	// Allow marking as needs_spare_parts without specifying parts for simplicity

	// Forward imagePaths to repo by extending its signature
	return s.repo.UpdateTechnicianStepWithImages(ticketID, stepID, technicianID, status, notes, sparePartsUsed, imagePaths)
}

// GetTechnicianStepProgress gets progress summary for a technician on a ticket
func (s *Service) GetTechnicianStepProgress(ticketID uint64, technicianID string) (map[string]interface{}, error) {
	return s.repo.GetTechnicianStepProgress(ticketID, technicianID)
}

// MarkTechnicianJobCompleted marks the entire technician job as completed
func (s *Service) MarkTechnicianJobCompleted(ticketID uint64) error {
	ok, err := s.repo.AreAllStepsCompleted(ticketID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("cannot mark job as completed: not all steps are done or marked not_applicable")
	}
	return s.repo.MarkTechnicianJobCompleted(ticketID)
}

// SetNetworkArchitecture sets FTTH/HTB/DISMANTLE for a ticket
func (s *Service) SetNetworkArchitecture(ticketID uint64, arch string) error {
	if arch != "FTTH" && arch != "HTB" && arch != "DISMANTLE" {
		return fmt.Errorf("invalid architecture: %s", arch)
	}
	return s.repo.SetNetworkArchitecture(ticketID, arch)
}

// UpsertTechnicianTeam validates and saves technician team
func (s *Service) UpsertTechnicianTeam(ticketID uint64, members []entities.TechnicianTeamMember) error {
	if len(members) == 0 {
		return fmt.Errorf("at least one technician must be assigned")
	}

	// Ensure distinct users per role
	seen := map[string]bool{}
	for _, m := range members {
		if m.Role != "senior" && m.Role != "junior" && m.Role != "helper" {
			return fmt.Errorf("invalid role: %s", m.Role)
		}
		if m.UserID == "" {
			return fmt.Errorf("user_id is required for role %s", m.Role)
		}
		if seen[m.UserID] {
			return fmt.Errorf("technician cannot be assigned to multiple roles in the same ticket")
		}
		seen[m.UserID] = true
		m.TicketID = ticketID
	}
	return s.repo.UpsertTechnicianTeam(ticketID, members)
}

// SaveSelfieStep saves a selfie photo as step 0 for a ticket
func (s *Service) SaveSelfieStep(ticketID uint64, technicianID string, imagePath string) error {
	// Get step with step_order = 0
	step, err := s.repo.GetStepByOrder(0)
	if err != nil {
		return fmt.Errorf("step 0 (selfie step) not found in technician_steps table. Please ensure a step with step_order=0 exists: %v", err)
	}

	// Convert single image path to JSON array
	imagePaths := []string{imagePath}
	imagePathsBytes, err := json.Marshal(imagePaths)
	if err != nil {
		return fmt.Errorf("failed to marshal image paths: %v", err)
	}
	imagePathsJSON := string(imagePathsBytes)

	// Save as done status with the image
	return s.repo.UpdateTechnicianStepWithImages(ticketID, step.ID, technicianID, "done", nil, nil, &imagePathsJSON)
}
