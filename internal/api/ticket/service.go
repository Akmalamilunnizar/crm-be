package ticketapi

import (
	"fmt"
	"log"
	"skripsi-be/internal/models/entities"
	"skripsi-be/internal/telegram"
)

type Service struct {
	repo            *Repo
	telegramService *telegram.Service
}

func NewService(r *Repo) *Service {
	return &Service{
		repo:            r,
		telegramService: telegram.NewService(),
	}
}

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

	// Send Telegram notification to CS channel
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", &input, "New Ticket Created"); err != nil {
		log.Printf("Failed to send Telegram notification: %v", err)
	}

	return &input, nil
}

// CreateTicketFromNetwatch creates a ticket from Netwatch monitoring
// DISABLED FOR TESTING - Netwatch integration is currently disabled
func (s *Service) CreateTicketFromNetwatch(input *entities.TroubleTicket) (*entities.TroubleTicket, error) {
	// Force initial state to match DB enum values
	input.Status = "ongoing" // Netwatch tickets start as ongoing

	// Look up CS role ID dynamically for assignment
	csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty")
	}
	input.CurrentAssignee = csRoleID

	if err := s.repo.Create(input); err != nil {
		return nil, err
	}
	return input, nil
}

func (s *Service) SendToNOC(id uint64, note string, imageFilename *string) (*entities.TroubleTicket, error) {
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
	t.CustomerNote = &note
	log.Printf("SendToNOC: Setting CustomerNote to: %s", note)
	if imageFilename != nil && *imageFilename != "" {
		t.ImgCS = imageFilename
		log.Printf("SendToNOC: Setting ImgCS to: %s", *imageFilename)
	} else {
		log.Printf("SendToNOC: No image provided")
	}
	log.Printf("SendToNOC: About to save ticket with CustomerNote=%v, ImgCS=%v", t.CustomerNote, t.ImgCS)
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}

	// Send Telegram notification to NOC channel
	if err := s.telegramService.SendTicketNotification("NOC", t, "Ticket Assigned to NOC"); err != nil {
		log.Printf("Failed to send Telegram notification to NOC: %v", err)
	}

	log.Printf("SendToNOC: Successfully saved ticket")
	return t, nil
}

// SendToCS moves the ticket back to Customer Service with an optional NOC note and image
func (s *Service) SendToCS(id uint64, note string, tType *string, imageFilename *string) (*entities.TroubleTicket, error) {
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

	// Save image filename if provided
	if imageFilename != nil && *imageFilename != "" {
		t.ImgNOC = imageFilename
	}

	if err := s.repo.Save(t); err != nil {
		return nil, err
	}

	// Send Telegram notification to CS channel
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", t, "Ticket Returned to Customer Service"); err != nil {
		log.Printf("Failed to send Telegram notification to CS: %v", err)
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

	// Send Telegram notification to CS channel
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", t, "Ticket Solved by NOC"); err != nil {
		log.Printf("Failed to send Telegram notification to CS: %v", err)
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

	// Send Telegram notification to CS channel
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", t, "Ticket Requires Physical Check"); err != nil {
		log.Printf("Failed to send Telegram notification to CS: %v", err)
	}

	return t, nil
}

func (s *Service) AssignTechnician(id uint64) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	// Don't assign to specific technician - leave AssignedTo as NULL so all technicians can see it
	t.AssignedTo = nil // Set to NULL so all technicians can see it
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

	// Send Telegram notification to Technician channel
	if err := s.telegramService.SendTicketNotification("TECHNICIAN", t, "Ticket Assigned to Technician"); err != nil {
		log.Printf("Failed to send Telegram notification to Technician: %v", err)
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

	// Send Telegram notification to CS channel
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", t, "Ticket Resolved by Technician"); err != nil {
		log.Printf("Failed to send Telegram notification to CS: %v", err)
	}

	return t, nil
}

// CSResolve allows Customer Service to resolve tickets with a customer note
func (s *Service) CSResolve(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	t.Status = "finished"
	// Keep assigned to CS since they're resolving it
	csRoleID, err := s.repo.RoleIDByName(string(entities.AssignCS))
	if err != nil {
		return nil, err
	}
	if csRoleID == "" {
		return nil, fmt.Errorf("cs role ID is empty")
	}
	t.CurrentAssignee = csRoleID
	t.CustomerNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

// AddTechnicianNote adds a technician note to a ticket without changing status
func (s *Service) AddTechnicianNote(id uint64, note string, imgTechBf *string, imgTechAf *string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}

	// Update the technician note field
	t.TechnicianNote = &note

	// Update image fields if provided
	if imgTechBf != nil {
		t.ImgTechBF = imgTechBf
	}
	if imgTechAf != nil {
		t.ImgTechAF = imgTechAf
	}

	if err := s.repo.Save(t); err != nil {
		return nil, err
	}

	// Send Telegram notification to CS channel about technician note
	if err := s.telegramService.SendTicketNotification("CUSTOMER SERVICE", t, "Technician Note Added"); err != nil {
		log.Printf("Failed to send Telegram notification to CS: %v", err)
	}

	return t, nil
}
