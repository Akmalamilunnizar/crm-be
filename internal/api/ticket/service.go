package ticketapi

import (
	"fmt"
	"log"
	"skripsi-be/internal/models/entities"
)

type Service struct{ repo *Repo }

func NewService(r *Repo) *Service { return &Service{r} }

func (s *Service) CreateCS(input entities.TroubleTicket) (*entities.TroubleTicket, error) {
	// Force initial state to match DB enum values
	input.Status = "unfinished"
	if input.CurrentAssignee == "" {
		// Look up CS role ID dynamically
		csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
		if err != nil {
			return nil, err
		}
		if csRoleID == "" {
			return nil, fmt.Errorf("cs role ID is empty")
		}
		input.CurrentAssignee = csRoleID
	}
	if err := s.repo.Create(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) SendToNOC(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.Status = "ongoing"
	// Look up NOC role ID dynamically
	nocRoleName := string(entities.AssignNOC)
	log.Printf("SendToNOC: Looking up role ID for name: '%s'", nocRoleName)
	nocRoleID, err := s.repo.RoleIDByName(nocRoleName)
	if err != nil {
		log.Printf("SendToNOC: Error looking up NOC role ID: %v", err)
		return nil, err
	}
	log.Printf("SendToNOC: Found NOC role ID: '%s'", nocRoleID)
	if nocRoleID == "" {
		return nil, fmt.Errorf("noc role ID is empty for role name: %s", nocRoleName)
	}
	t.CurrentAssignee = nocRoleID
	t.NOCNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

// SendToCS moves the ticket back to Customer Service with an optional NOC note
func (s *Service) SendToCS(id uint64, note string, tType *string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	// keep status ongoing while reassigning
	t.Status = "ongoing"
	// If NOC provided a diagnosed trouble type, persist it
	if tType != nil && *tType != "" {
		t.Type = tType
	}
	// Look up CS role ID dynamically
	csRoleName := string(entities.AssignCS)
	log.Printf("SendToCS: Looking up role ID for name: '%s'", csRoleName)
	csRoleID, err := s.repo.RoleIDByName(csRoleName)
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty for role name: %s", csRoleName)
	}
	t.CurrentAssignee = csRoleID
	// reuse NOC note field for context when NOC sends back to CS
	t.NOCNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) NOCSolved(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.Status = "finished"
	// Look up CS role ID dynamically
	csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty")
	}
	t.CurrentAssignee = csRoleID
	t.NOCNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) NOCPhysical(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.Status = "ongoing"
	// Look up CS role ID dynamically
	csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty")
	}
	t.CurrentAssignee = csRoleID
	t.NOCNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) AssignTechnician(id uint64, techUserID string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.AssignedTo = techUserID
	// Look up Technician role ID dynamically
	techRoleID, err := s.repo.RoleIDByName(string(entities.AssignTech))
	if err != nil {
		return nil, err
	}
	if techRoleID == "" {
		return nil, fmt.Errorf("technician role ID is empty")
	}
	t.CurrentAssignee = techRoleID
	if t.Status == "unfinished" {
		t.Status = "ongoing"
	}
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) TechnicianResolve(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.Status = "finished"
	// Look up CS role ID dynamically
	csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty")
	}
	t.CurrentAssignee = csRoleID
	t.TechnicianNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}
