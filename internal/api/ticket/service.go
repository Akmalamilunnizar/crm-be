package ticketapi

import "skripsi-be/internal/models/entities"

type Service struct{ repo *Repo }

func NewService(r *Repo) *Service { return &Service{r} }

func (s *Service) CreateCS(input entities.TroubleTicket) (*entities.TroubleTicket, error) {
	// Force initial state to match DB enum values
	input.Status = "unfinished"
	if input.CurrentAssignee == "" {
		input.CurrentAssignee = string(entities.AssignCS)
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
	t.CurrentAssignee = string(entities.AssignNOC)
	t.NOCNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

// SendToCS moves the ticket back to Customer Service with an optional NOC note
func (s *Service) SendToCS(id uint64, note string) (*entities.TroubleTicket, error) {
	t, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}
	// keep status ongoing while reassigning
	t.Status = "ongoing"
	t.CurrentAssignee = string(entities.AssignCS)
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
	t.CurrentAssignee = string(entities.AssignCS)
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
	t.CurrentAssignee = string(entities.AssignCS)
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
	t.CurrentAssignee = string(entities.AssignTech)
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
	t.CurrentAssignee = string(entities.AssignCS)
	t.TechnicianNote = &note
	if err := s.repo.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}
