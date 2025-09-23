package ticketapi

import (
	"fmt"
	"log"
	"skripsi-be/internal/models/entities"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repo struct{ DB *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db} }

// TicketWithAssignee includes assignee name and customer name for display
type TicketWithAssignee struct {
	entities.TroubleTicket
	AssigneeName        string   `json:"assignee_name" gorm:"column:assignee_name"`
	CurrentAssigneeName string   `json:"current_assignee_name" gorm:"column:current_assignee_name"`
	TypeName            string   `json:"type_name" gorm:"column:type_name"`
	CustomerName        string   `json:"customer_name" gorm:"column:customer_name"`
	GPSLat              *float64 `json:"gps_lat" gorm:"column:gps_lat"`
	GPSLng              *float64 `json:"gps_lng" gorm:"column:gps_lng"`
	CustomerAddress     *string  `json:"customer_address" gorm:"column:customer_address"`
	CustomerPhone       *string  `json:"customer_phone" gorm:"column:customer_phone"`
}

func (r *Repo) ListAll() ([]TicketWithAssignee, error) {
	// Debug: Check what's in trouble_type table
	var troubleTypes []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	r.DB.Table("trouble_type").Select("id, name").Scan(&troubleTypes)
	log.Printf("Trouble types in database: %+v", troubleTypes)

	var items []TicketWithAssignee
	err := r.DB.Table("trouble_tickets t").
		Select("t.*, r1.name as assignee_name, r2.name as current_assignee_name, tt.name as type_name, c.name as customer_name, c.latitude as gps_lat, c.longitude as gps_lng, c.address as customer_address, c.phone as customer_phone").
		Joins("LEFT JOIN users r1 ON r1.id = t.assigned_to").
		Joins("LEFT JOIN roles r2 ON r2.id = t.current_assignee_role").
		Joins("LEFT JOIN trouble_type tt ON tt.id = t.type").
		Joins("LEFT JOIN customer c ON c.id = t.customer_id").
		Order("t.created_at DESC").
		Scan(&items).Error

	// Debug logging
	for i, item := range items {
		assignedTo := "NULL"
		if item.AssignedTo != nil {
			assignedTo = *item.AssignedTo
		}
		log.Printf("Ticket %d: ID=%d, CustomerID=%s, CustomerName='%s', AssignedTo=%s, AssigneeName='%s', CurrentAssigneeRole='%s', CurrentAssigneeName='%s', Type='%s', TypeName='%s'",
			i+1, item.ID, item.CustomerID, item.CustomerName, assignedTo, item.AssigneeName, item.CurrentAssignee, item.CurrentAssigneeName, item.Type, item.TypeName)
	}

	return items, err
}

// Debug method to check roles and assigned_to values
func (r *Repo) DebugRolesAndAssignments() {
	// Check what roles exist
	var roles []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	r.DB.Table("roles").Select("id, name").Scan(&roles)
	log.Printf("Roles in database: %+v", roles)

	// Check what assigned_to values exist in tickets
	var assignments []struct {
		ID         uint64 `json:"id"`
		AssignedTo string `json:"assigned_to"`
	}
	r.DB.Table("trouble_tickets").Select("id, assigned_to").Scan(&assignments)
	log.Printf("Ticket assignments: %+v", assignments)
}

func (r *Repo) ByID(id uint64) (*entities.TroubleTicket, error) {
	var t entities.TroubleTicket
	return &t, r.DB.First(&t, "id = ?", id).Error
}

func (r *Repo) Create(t *entities.TroubleTicket) error { return r.DB.Create(t).Error }
func (r *Repo) Save(t *entities.TroubleTicket) error   { return r.DB.Save(t).Error }

// trouble_type CRUD
func (r *Repo) CreateTroubleType(t entities.TroubleTypeRow) error {
	return r.DB.Table("trouble_type").Create(&t).Error
}

// reporting

type ReportSlice struct {
	Type string `json:"type"`
	Cnt  int64  `json:"count"`
}

func (r *Repo) CountByType(startDate, endDate string) ([]ReportSlice, error) {
	var rows []ReportSlice
	query := r.DB.Table("trouble_tickets")

	// Add date filtering if dates are provided
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Select("type, COUNT(*) as cnt").
		Group("type").Scan(&rows).Error
	return rows, err
}

// Lookup helpers
func (r *Repo) AllTroubleTypes() ([]entities.TroubleTypeRow, error) {
	var rows []entities.TroubleTypeRow
	return rows, r.DB.Table("trouble_type").Order("name asc").Find(&rows).Error
}

// Technician Team operations
func (r *Repo) ReplaceTeamMembers(ticketID uint64, members []entities.TechnicianTeamMember) error {
	tx := r.DB.Begin()
	if err := tx.Table((entities.TechnicianTeamMember{}).TableName()).Where("ticket_id = ?", ticketID).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}
	for i := range members {
		members[i].TicketID = ticketID
		if err := tx.Create(&members[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *Repo) ListTeamMembers(ticketID uint64) ([]entities.TechnicianTeamMember, error) {
	var rows []entities.TechnicianTeamMember
	err := r.DB.Table((entities.TechnicianTeamMember{}).TableName()).Where("ticket_id = ?", ticketID).Find(&rows).Error
	return rows, err
}

// Ticket Steps operations
func (r *Repo) AddTicketStep(step *entities.TicketStep) error {
	return r.DB.Table((entities.TicketStep{}).TableName()).Create(step).Error
}

func (r *Repo) ListTicketSteps(ticketID uint64) ([]entities.TicketStep, error) {
	var rows []entities.TicketStep
	err := r.DB.Table((entities.TicketStep{}).TableName()).Where("ticket_id = ?", ticketID).Order("step_order asc").Find(&rows).Error
	return rows, err
}

func (r *Repo) CountTicketSteps(ticketID uint64) (int64, error) {
	var cnt int64
	err := r.DB.Table((entities.TicketStep{}).TableName()).Where("ticket_id = ?", ticketID).Count(&cnt).Error
	return cnt, err
}

// Step images
func (r *Repo) AddStepImage(img *entities.TicketStepImage) error {
	return r.DB.Table((entities.TicketStepImage{}).TableName()).Create(img).Error
}

// Hotspot aggregation by customer geolocation
type Hotspot struct {
	GPSLat *float64 `json:"gps_lat"`
	GPSLng *float64 `json:"gps_lng"`
	Cnt    int64    `json:"count"`
}

func (r *Repo) HotLocations() ([]Hotspot, error) {
	var rows []Hotspot
	err := r.DB.Table("trouble_tickets t").
		Select("c.latitude as gps_lat, c.longitude as gps_lng, COUNT(*) as cnt").
		Joins("JOIN customer c ON c.id = t.customer_id").
		Where("c.latitude IS NOT NULL AND c.longitude IS NOT NULL").
		Group("c.latitude, c.longitude").
		Order("cnt DESC").
		Limit(50).
		Scan(&rows).Error
	return rows, err
}

// RoleIDByName finds role id by roles.name with fallback between space/underscore variants
func (r *Repo) RoleIDByName(name string) (string, error) {
	log.Printf("RoleIDByName: Looking for role with name: '%s'", name)
	var row struct{ ID string }
	if err := r.DB.Table("roles").Select("id").Where("name = ?", name).Scan(&row).Error; err != nil {
		log.Printf("RoleIDByName: Error querying role '%s': %v", name, err)
		return "", err
	}
	if row.ID == "" {
		// Try space/underscore normalized counterpart
		alt := name
		if strings.Contains(name, " ") {
			alt = strings.ReplaceAll(name, " ", "_")
		} else if strings.Contains(name, "_") {
			alt = strings.ReplaceAll(name, "_", " ")
		}
		if alt != name {
			log.Printf("RoleIDByName: Not found, trying alternative name: '%s'", alt)
			if err := r.DB.Table("roles").Select("id").Where("name = ?", alt).Scan(&row).Error; err != nil {
				log.Printf("RoleIDByName: Error querying role alt '%s': %v", alt, err)
				return "", err
			}
		}
	}
	if row.ID == "" {
		err := fmt.Errorf("role '%s' not found (after fallback)", name)
		log.Printf("RoleIDByName: %v", err)
		return "", err
	}
	log.Printf("RoleIDByName: Found role ID '%s' for name '%s'", row.ID, name)
	return row.ID, nil
}

// UpdatesSinceForRole returns tickets updated recently and relevant to the caller's role
func (r *Repo) UpdatesSinceForRole(since time.Time, normalizedRole string, userID string) ([]TicketWithAssignee, error) {
	var items []TicketWithAssignee

	q := r.DB.Table("trouble_tickets t").
		Select("t.*, r1.name as assignee_name, r2.name as current_assignee_name").
		Joins("LEFT JOIN users r1 ON r1.id = t.assigned_to").
		Joins("LEFT JOIN roles r2 ON r2.id = t.current_assignee_role").
		Where("(t.updated_at > ? OR t.created_at > ?)", since, since)

	switch normalizedRole {
	case "ADMIN":
		// admins see everything
	case "NOC", "CUSTOMER_SERVICE", "TECHNICIAN":
		roleName := normalizedRole
		if normalizedRole == "CUSTOMER_SERVICE" {
			roleName = "CUSTOMER SERVICE"
		}
		roleID, err := r.RoleIDByName(roleName)
		if err != nil {
			return items, err
		}
		if normalizedRole == "TECHNICIAN" {
			q = q.Where("t.assigned_to = ? OR t.current_assignee_role = ?", userID, roleID)
		} else {
			q = q.Where("t.current_assignee_role = ?", roleID)
		}
	default:
		return items, nil
	}

	if err := q.Order("t.updated_at DESC").Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MarkTechnicianCompleted marks a ticket as completed by technician and reassigns to CS
func (r *Repo) MarkTechnicianCompleted(ticketID uint64, technicianUserID string) (*entities.TroubleTicket, error) {
	// Get CS role ID
	csRoleID, err := r.RoleIDByName("CUSTOMER SERVICE")
	if err != nil {
		return nil, fmt.Errorf("failed to get CS role ID: %v", err)
	}

	// Update ticket: mark as completed and reassign to CS
	var ticket entities.TroubleTicket
	if err := r.DB.Model(&ticket).Where("id = ? AND assigned_to = ?", ticketID, technicianUserID).
		Updates(map[string]interface{}{
			"technician_completed":  true,
			"current_assignee_role": csRoleID,
			"assigned_to":           nil, // Remove technician assignment
			"updated_at":            time.Now(),
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to mark technician completed: %v", err)
	}

	// Get updated ticket
	if err := r.DB.First(&ticket, ticketID).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated ticket: %v", err)
	}

	return &ticket, nil
}

// (removed duplicate SetNetworkArchitecture; see technician_repository.go)

// ValidateTechnicianAssignment validates that a technician is not assigned to multiple roles in the same ticket
func (r *Repo) ValidateTechnicianAssignment(ticketID uint64, technicianUserID string) error {
	var count int64
	if err := r.DB.Model(&entities.TechnicianTeamMember{}).
		Where("ticket_id = ? AND user_id = ?", ticketID, technicianUserID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate technician assignment: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("technician %s is already assigned to this ticket", technicianUserID)
	}

	return nil
}

// GetTechnicianTeamMembers gets all team members for a ticket
func (r *Repo) GetTechnicianTeamMembers(ticketID uint64) ([]entities.TechnicianTeamMember, error) {
	var members []entities.TechnicianTeamMember
	if err := r.DB.Where("ticket_id = ?", ticketID).Find(&members).Error; err != nil {
		return nil, fmt.Errorf("failed to get team members: %v", err)
	}
	return members, nil
}
