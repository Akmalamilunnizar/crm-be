package ticketapi

import (
	"fmt"
	"skripsi-be/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

// GetTechnicianSteps gets all predefined technician steps
func (r *Repo) GetTechnicianSteps() ([]entities.TechnicianStep, error) {
	var steps []entities.TechnicianStep
	if err := r.DB.Where("is_active = ?", true).Order("step_order ASC").Find(&steps).Error; err != nil {
		return nil, fmt.Errorf("failed to get technician steps: %v", err)
	}
	return steps, nil
}

// GetSpareParts gets all available spare parts
func (r *Repo) GetSpareParts() ([]entities.SparePart, error) {
	var parts []entities.SparePart
	if err := r.DB.Where("is_active = ?", true).Find(&parts).Error; err != nil {
		return nil, fmt.Errorf("failed to get spare parts: %v", err)
	}
	return parts, nil
}

// GetTicketTechnicianSteps gets technician progress for a specific ticket
func (r *Repo) GetTicketTechnicianSteps(ticketID uint64, technicianID string) ([]entities.TicketTechnicianStep, error) {
	var steps []entities.TicketTechnicianStep
	if err := r.DB.Preload("Step").Where("ticket_id = ? AND technician_id = ?", ticketID, technicianID).Find(&steps).Error; err != nil {
		return nil, fmt.Errorf("failed to get ticket technician steps: %v", err)
	}
	return steps, nil
}

// UpdateTechnicianStep updates a technician's progress on a specific step
func (r *Repo) UpdateTechnicianStepWithImages(ticketID uint64, stepID uint64, technicianID string, status string, notes *string, sparePartsUsed *string, imagePaths *string) error {
	now := time.Now()

	// Check if step already exists
	var existingStep entities.TicketTechnicianStep
	err := r.DB.Where("ticket_id = ? AND step_id = ? AND technician_id = ?", ticketID, stepID, technicianID).First(&existingStep).Error

	if err == gorm.ErrRecordNotFound {
		// Create new step
		step := entities.TicketTechnicianStep{
			TicketID:       ticketID,
			StepID:         stepID,
			TechnicianID:   technicianID,
			Status:         status,
			Notes:          notes,
			SparePartsUsed: sparePartsUsed,
			ImagePaths:     imagePaths,
			CreatedAt:      &now,
			UpdatedAt:      &now,
		}

		if status == "done" {
			step.CompletedAt = &now
		}

		if err := r.DB.Create(&step).Error; err != nil {
			return fmt.Errorf("failed to create technician step: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing step: %v", err)
	} else {
		// Update existing step
		updates := map[string]interface{}{
			"status":     status,
			"notes":      notes,
			"updated_at": now,
		}

		if sparePartsUsed != nil {
			updates["spare_parts_used"] = sparePartsUsed
		}
		if imagePaths != nil {
			updates["image_paths"] = imagePaths
		}

		if status == "done" && existingStep.Status != "done" {
			updates["completed_at"] = now
		}

		if err := r.DB.Model(&existingStep).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update technician step: %v", err)
		}
	}

	return nil
}

// GetTicketByID fetches a ticket by id
func (r *Repo) GetTicketByID(ticketID uint64) (*entities.TroubleTicket, error) {
	var t entities.TroubleTicket
	if err := r.DB.Where("id = ?", ticketID).First(&t).Error; err != nil {
		return nil, fmt.Errorf("failed to get ticket: %v", err)
	}
	return &t, nil
}

// GetStepByID returns a predefined step
func (r *Repo) GetStepByID(stepID uint64) (*entities.TechnicianStep, error) {
	var s entities.TechnicianStep
	if err := r.DB.Where("id = ?", stepID).First(&s).Error; err != nil {
		return nil, fmt.Errorf("failed to get step: %v", err)
	}
	return &s, nil
}

// HasIncompletePreviousSteps checks if there are previous steps not completed
func (r *Repo) HasIncompletePreviousSteps(ticketID uint64, technicianID string, currentOrder int) (bool, error) {
	var count int64
	// Any step with order < currentOrder must be done or not_applicable
	// We join to technician_steps to compare order
	if err := r.DB.Model(&entities.TicketTechnicianStep{}).
		Joins("JOIN technician_steps ts ON ts.id = ticket_technician_steps.step_id").
		Where("ticket_technician_steps.ticket_id = ? AND ticket_technician_steps.technician_id = ? AND ts.step_order < ? AND ticket_technician_steps.status NOT IN (?)",
			ticketID, technicianID, currentOrder, []string{"done", "not_applicable", "needs_spare_parts"}).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check previous steps: %v", err)
	}
	return count > 0, nil
}

// SetNetworkArchitecture updates ticket network architecture
func (r *Repo) SetNetworkArchitecture(ticketID uint64, arch string) error {
	if err := r.DB.Model(&entities.TroubleTicket{}).Where("id = ?", ticketID).Update("network_architecture", arch).Error; err != nil {
		return fmt.Errorf("failed to set network architecture: %v", err)
	}
	return nil
}

// UpsertTechnicianTeam replaces team roles for a ticket
func (r *Repo) UpsertTechnicianTeam(ticketID uint64, members []entities.TechnicianTeamMember) error {
	tx := r.DB.Begin()
	if err := tx.Where("ticket_id = ?", ticketID).Delete(&entities.TechnicianTeamMember{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear team: %v", err)
	}
	if len(members) > 0 {
		if err := tx.Create(&members).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save team: %v", err)
		}
	}
	return tx.Commit().Error
}

// GetTechnicianStepProgress gets progress summary for a technician on a ticket
func (r *Repo) GetTechnicianStepProgress(ticketID uint64, technicianID string) (map[string]interface{}, error) {
	var totalSteps int64
	var completedSteps int64
	var pendingSteps int64
	var needsSpareParts int64

	// Count total steps
	if err := r.DB.Model(&entities.TechnicianStep{}).Where("is_active = ?", true).Count(&totalSteps).Error; err != nil {
		return nil, fmt.Errorf("failed to count total steps: %v", err)
	}

	// Count completed steps
	if err := r.DB.Model(&entities.TicketTechnicianStep{}).
		Where("ticket_id = ? AND technician_id = ? AND status = ?", ticketID, technicianID, "done").
		Count(&completedSteps).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed steps: %v", err)
	}

	// Count pending steps
	if err := r.DB.Model(&entities.TicketTechnicianStep{}).
		Where("ticket_id = ? AND technician_id = ? AND status = ?", ticketID, technicianID, "pending").
		Count(&pendingSteps).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending steps: %v", err)
	}

	// Count steps needing spare parts
	if err := r.DB.Model(&entities.TicketTechnicianStep{}).
		Where("ticket_id = ? AND technician_id = ? AND status = ?", ticketID, technicianID, "needs_spare_parts").
		Count(&needsSpareParts).Error; err != nil {
		return nil, fmt.Errorf("failed to count steps needing spare parts: %v", err)
	}

	progress := float64(completedSteps) / float64(totalSteps) * 100

	return map[string]interface{}{
		"total_steps":         totalSteps,
		"completed_steps":     completedSteps,
		"pending_steps":       pendingSteps,
		"needs_spare_parts":   needsSpareParts,
		"progress_percentage": progress,
		"is_completed":        completedSteps == totalSteps,
	}, nil
}

// MarkTechnicianJobCompleted marks the entire technician job as completed
func (r *Repo) MarkTechnicianJobCompleted(ticketID uint64) error {
	// Set technician_completed and mark ticket finished
	updates := map[string]interface{}{
		"technician_completed": true,
		"status":               "finished",
		"updated_at":           time.Now(),
	}
	if err := r.DB.Model(&entities.TroubleTicket{}).Where("id = ?", ticketID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to mark technician job as completed: %v", err)
	}
	return nil
}

// AreAllStepsCompleted checks if all predefined steps have at least one record done or fix (needs_spare_parts) or not_applicable
func (r *Repo) AreAllStepsCompleted(ticketID uint64) (bool, error) {
	// Count active steps
	var total int64
	if err := r.DB.Model(&entities.TechnicianStep{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return false, fmt.Errorf("failed to count steps: %v", err)
	}
	if total == 0 {
		return false, nil
	}
	// Count how many step IDs have a done/not_applicable record for this ticket
	type row struct{ StepID uint64 }
	var rows []row
	if err := r.DB.Model(&entities.TicketTechnicianStep{}).
		Select("DISTINCT step_id").
		Where("ticket_id = ? AND status IN (?)", ticketID, []string{"done", "not_applicable", "needs_spare_parts"}).
		Find(&rows).Error; err != nil {
		return false, fmt.Errorf("failed to query completion: %v", err)
	}
	return int64(len(rows)) == total, nil
}
