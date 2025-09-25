package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"skripsi-be/internal/models/entities"
	"gorm.io/gorm"
)

// AccumulationService handles automatic detection and grouping of similar troubles
type AccumulationService struct {
	db *gorm.DB
}

// NewAccumulationService creates a new accumulation service
func NewAccumulationService(db *gorm.DB) *AccumulationService {
	return &AccumulationService{db: db}
}

// SimilarTrouble represents a trouble that might be similar to others
type SimilarTrouble struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Accumulation int      `json:"accumulation"`
}

// DetectSimilarTroubles finds tickets that might be related to the same problem
func (s *AccumulationService) DetectSimilarTroubles(ticketID uint64, timeWindowHours int) ([]SimilarTrouble, error) {
	var currentTicket SimilarTrouble
	
	// Get the current ticket details
	err := s.db.Table("trouble_tickets").
		Select("id, title, type, description, created_at, accumulation").
		Where("id = ?", ticketID).
		First(&currentTicket).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get current ticket: %v", err)
	}
	
	// Calculate time window
	timeWindow := time.Duration(timeWindowHours) * time.Hour
	startTime := currentTicket.CreatedAt.Add(-timeWindow)
	endTime := currentTicket.CreatedAt.Add(timeWindow)
	
	var similarTickets []SimilarTrouble
	
	// Find tickets with similar characteristics within the time window
	query := s.db.Table("trouble_tickets").
		Select("id, title, type, description, created_at, accumulation").
		Where("id != ?", ticketID).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Where("status IN (?)", []string{"unfinished", "ongoing"})
	
	// Add similarity conditions
	similarityConditions := s.buildSimilarityConditions(currentTicket)
	if len(similarityConditions) > 0 {
		query = query.Where(strings.Join(similarityConditions, " OR "))
	}
	
	err = query.Find(&similarTickets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find similar tickets: %v", err)
	}
	
	return similarTickets, nil
}

// buildSimilarityConditions creates SQL conditions for finding similar tickets
func (s *AccumulationService) buildSimilarityConditions(ticket SimilarTrouble) []string {
	var conditions []string
	
	// Same trouble type
	if ticket.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = '%s'", ticket.Type))
	}
	
	// Similar title keywords (basic keyword matching)
	titleKeywords := s.extractKeywords(ticket.Title)
	for _, keyword := range titleKeywords {
		if len(keyword) > 3 { // Only consider meaningful keywords
			conditions = append(conditions, fmt.Sprintf("title LIKE '%%%s%%'", keyword))
		}
	}
	
	// Similar description keywords
	descKeywords := s.extractKeywords(ticket.Description)
	for _, keyword := range descKeywords {
		if len(keyword) > 3 {
			conditions = append(conditions, fmt.Sprintf("description LIKE '%%%s%%'", keyword))
		}
	}
	
	return conditions
}

// extractKeywords extracts meaningful keywords from text
func (s *AccumulationService) extractKeywords(text string) []string {
	if text == "" {
		return []string{}
	}
	
	// Convert to lowercase and split by common delimiters
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(c rune) bool {
		return c == ' ' || c == ',' || c == '.' || c == '!' || c == '?' || c == '-'
	})
	
	var keywords []string
	// Filter out common words and short words
	commonWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true, "should": true,
		"tidak": true, "ada": true, "yang": true, "dan": true, "atau": true, "dengan": true,
		"untuk": true, "pada": true, "di": true, "ke": true,
	}
	
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 3 && !commonWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// UpdateAccumulation updates the accumulation count for a group of similar tickets
func (s *AccumulationService) UpdateAccumulation(ticketIDs []uint64, accumulation int) error {
	if len(ticketIDs) == 0 {
		return fmt.Errorf("no ticket IDs provided")
	}
	
	// Update all tickets with the same accumulation count
	err := s.db.Model(&struct {
		ID uint64 `gorm:"primaryKey"`
	}{}).
		Table("trouble_tickets").
		Where("id IN (?)", ticketIDs).
		Update("accumulation", accumulation).Error
	
	if err != nil {
		return fmt.Errorf("failed to update accumulation: %v", err)
	}
	
	log.Printf("Updated accumulation to %d for %d tickets: %v", accumulation, len(ticketIDs), ticketIDs)
	return nil
}

// AutoDetectAndGroup automatically detects and groups similar troubles
func (s *AccumulationService) AutoDetectAndGroup() error {
	log.Println("Starting automatic trouble detection and grouping...")
	
	// Get all unfinished and ongoing tickets from the last 24 hours
	var tickets []SimilarTrouble
	err := s.db.Table("trouble_tickets").
		Select("id, title, type, description, created_at, accumulation").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Where("created_at >= ?", time.Now().Add(-24*time.Hour)).
		Order("created_at ASC").
		Find(&tickets).Error
	
	if err != nil {
		return fmt.Errorf("failed to get recent tickets: %v", err)
	}
	
	log.Printf("Found %d tickets to analyze", len(tickets))
	
	// For massive scale operations (500+ tickets), use batch processing
	if len(tickets) > 500 {
		log.Println("Large scale operation detected - using optimized batch processing")
		return s.batchProcessLargeScaleSimple(tickets)
	}
	
	// Group tickets by similarity
	groups := s.groupSimilarTickets(tickets)
	
	// Update accumulation for each group
	for _, group := range groups {
		if len(group) > 1 {
			err := s.UpdateAccumulation(group, len(group))
			if err != nil {
				log.Printf("Failed to update accumulation for group %v: %v", group, err)
			}
		}
	}
	
	log.Printf("Automatic grouping completed. Processed %d groups", len(groups))
	return nil
}

// groupSimilarTickets groups tickets that are likely related
func (s *AccumulationService) groupSimilarTickets(tickets []SimilarTrouble) [][]uint64 {
	var groups [][]uint64
	processed := make(map[uint64]bool)
	
	for i, ticket := range tickets {
		if processed[ticket.ID] {
			continue
		}
		
		group := []uint64{ticket.ID}
		processed[ticket.ID] = true
		
		// Find similar tickets
		for j := i + 1; j < len(tickets); j++ {
			otherTicket := tickets[j]
			if processed[otherTicket.ID] {
				continue
			}
			
			if s.areTicketsSimilar(ticket, otherTicket) {
				group = append(group, otherTicket.ID)
				processed[otherTicket.ID] = true
			}
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

// areTicketsSimilar determines if two tickets are similar enough to be grouped
func (s *AccumulationService) areTicketsSimilar(ticket1, ticket2 SimilarTrouble) bool {
	// Same type is a strong indicator
	if ticket1.Type == ticket2.Type && ticket1.Type != "" {
		return true
	}
	
	// Similar titles (basic keyword overlap)
	title1Keywords := s.extractKeywords(ticket1.Title)
	title2Keywords := s.extractKeywords(ticket2.Title)
	
	commonKeywords := 0
	for _, kw1 := range title1Keywords {
		for _, kw2 := range title2Keywords {
			if kw1 == kw2 {
				commonKeywords++
				break
			}
		}
	}
	
	// If more than 2 common keywords, consider similar
	if commonKeywords >= 2 {
		return true
	}
	
	// Similar descriptions
	desc1Keywords := s.extractKeywords(ticket1.Description)
	desc2Keywords := s.extractKeywords(ticket2.Description)
	
	commonDescKeywords := 0
	for _, kw1 := range desc1Keywords {
		for _, kw2 := range desc2Keywords {
			if kw1 == kw2 {
				commonDescKeywords++
				break
			}
		}
	}
	
	// If more than 3 common description keywords, consider similar
	if commonDescKeywords >= 3 {
		return true
	}
	
	return false
}

// GetAccumulationStats returns statistics about ticket accumulation
func (s *AccumulationService) GetAccumulationStats() (map[string]interface{}, error) {
	var stats struct {
		TotalTickets        int `json:"total_tickets"`
		SingleCustomer      int `json:"single_customer"`
		MultipleCustomers   int `json:"multiple_customers"`
		TotalCustomersAffected int `json:"total_customers_affected"`
		MaxAccumulation     int `json:"max_accumulation"`
	}
	
	// Get basic counts
	err := s.db.Table("trouble_tickets").
		Select("COUNT(*) as total_tickets").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Scan(&stats.TotalTickets).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get total tickets: %v", err)
	}
	
	// Get single customer tickets
	err = s.db.Table("trouble_tickets").
		Select("COUNT(*) as single_customer").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Where("accumulation = 1").
		Scan(&stats.SingleCustomer).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get single customer count: %v", err)
	}
	
	// Get multiple customer tickets
	err = s.db.Table("trouble_tickets").
		Select("COUNT(*) as multiple_customers").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Where("accumulation > 1").
		Scan(&stats.MultipleCustomers).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple customer count: %v", err)
	}
	
	// Get total customers affected
	err = s.db.Table("trouble_tickets").
		Select("SUM(accumulation) as total_customers_affected").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Scan(&stats.TotalCustomersAffected).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get total customers affected: %v", err)
	}
	
	// Get max accumulation
	err = s.db.Table("trouble_tickets").
		Select("MAX(accumulation) as max_accumulation").
		Where("status IN (?)", []string{"unfinished", "ongoing"}).
		Scan(&stats.MaxAccumulation).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get max accumulation: %v", err)
	}
	
	return map[string]interface{}{
		"total_tickets": stats.TotalTickets,
		"single_customer": stats.SingleCustomer,
		"multiple_customers": stats.MultipleCustomers,
		"total_customers_affected": stats.TotalCustomersAffected,
		"max_accumulation": stats.MaxAccumulation,
	}, nil
}

// batchProcessLargeScale handles massive scale operations (500+ tickets) efficiently
func (s *AccumulationService) batchProcessLargeScale(tickets []entities.TroubleTicket) error {
	log.Printf("Processing large scale operation with %d tickets", len(tickets))
	
	// Group tickets by type first for efficiency
	typeGroups := make(map[string][]entities.TroubleTicket)
	for _, ticket := range tickets {
		ticketType := "other"
		if ticket.Type != nil {
			ticketType = *ticket.Type
		}
		typeGroups[ticketType] = append(typeGroups[ticketType], ticket)
	}
	
	// Process each type group
	for ticketType, typeTickets := range typeGroups {
		if len(typeTickets) < 2 {
			continue
		}
		
		log.Printf("Processing %d tickets of type: %s", len(typeTickets), ticketType)
		
		// For large groups, use simpler similarity detection
		groups := s.groupSimilarTicketsSimple(typeTickets)
		
		// Update accumulation for each group
		for _, group := range groups {
			if len(group) > 1 {
				accumulation := len(group)
				ticketIDs := make([]uint64, len(group))
				for i, ticket := range group {
					ticketIDs[i] = ticket.ID
				}
				
				err := s.UpdateAccumulation(ticketIDs, accumulation)
				if err != nil {
					log.Printf("Failed to update accumulation for group: %v", err)
					continue
				}
				
				log.Printf("Updated group of %d similar tickets with accumulation %d", len(group), accumulation)
			}
		}
	}
	
	log.Println("Large scale processing completed")
	return nil
}

// groupSimilarTicketsSimple is a simplified version for large scale operations
func (s *AccumulationService) groupSimilarTicketsSimple(tickets []entities.TroubleTicket) [][]entities.TroubleTicket {
	var groups [][]entities.TroubleTicket
	processed := make(map[uint64]bool)
	
	for i, ticket := range tickets {
		if processed[ticket.ID] {
			continue
		}
		
		group := []entities.TroubleTicket{ticket}
		processed[ticket.ID] = true
		
		// Find similar tickets within the same time window (1 hour)
		timeWindow := ticket.CreatedAt.Add(-1 * time.Hour)
		
		for j := i + 1; j < len(tickets); j++ {
			otherTicket := tickets[j]
			if processed[otherTicket.ID] {
				continue
			}
			
			// Check if within time window
			if otherTicket.CreatedAt.Before(timeWindow) {
				continue
			}
			
			// Simple similarity check - same type and similar title keywords
			if s.areTicketsSimilarSimple(ticket, otherTicket) {
				group = append(group, otherTicket)
				processed[otherTicket.ID] = true
			}
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

// areTicketsSimilarSimple is a simplified similarity check for large scale operations
func (s *AccumulationService) areTicketsSimilarSimple(ticket1, ticket2 entities.TroubleTicket) bool {
	// Same type is a strong indicator
	if ticket1.Type != nil && ticket2.Type != nil && *ticket1.Type == *ticket2.Type {
		return true
	}
	
	// Simple keyword overlap in titles
	title1Keywords := s.extractKeywords(ticket1.Title)
	title2Keywords := s.extractKeywords(ticket2.Title)
	
	commonKeywords := 0
	for _, kw1 := range title1Keywords {
		for _, kw2 := range title2Keywords {
			if kw1 == kw2 {
				commonKeywords++
				break
			}
		}
	}
	
	// If more than 1 common keyword, consider similar
	return commonKeywords >= 1
}

// batchProcessLargeScaleSimple handles massive scale operations for SimilarTrouble structs
func (s *AccumulationService) batchProcessLargeScaleSimple(tickets []SimilarTrouble) error {
	log.Printf("Processing large scale operation with %d tickets", len(tickets))
	
	// Group tickets by type first for efficiency
	typeGroups := make(map[string][]SimilarTrouble)
	for _, ticket := range tickets {
		ticketType := "other"
		if ticket.Type != "" {
			ticketType = ticket.Type
		}
		typeGroups[ticketType] = append(typeGroups[ticketType], ticket)
	}
	
	// Process each type group
	for ticketType, typeTickets := range typeGroups {
		if len(typeTickets) < 2 {
			continue
		}
		
		log.Printf("Processing %d tickets of type: %s", len(typeTickets), ticketType)
		
		// For large groups, use simpler similarity detection
		groups := s.groupSimilarTicketsSimpleSimilar(typeTickets)
		
		// Update accumulation for each group
		for _, group := range groups {
			if len(group) > 1 {
				err := s.UpdateAccumulation(group, len(group))
				if err != nil {
					log.Printf("Failed to update accumulation for group: %v", err)
					continue
				}
				
				log.Printf("Updated group of %d similar tickets with accumulation %d", len(group), len(group))
			}
		}
	}
	
	log.Println("Large scale processing completed")
	return nil
}

// groupSimilarTicketsSimpleSimilar is a simplified version for SimilarTrouble structs
func (s *AccumulationService) groupSimilarTicketsSimpleSimilar(tickets []SimilarTrouble) [][]uint64 {
	var groups [][]uint64
	processed := make(map[uint64]bool)
	
	for i, ticket := range tickets {
		if processed[ticket.ID] {
			continue
		}
		
		group := []uint64{ticket.ID}
		processed[ticket.ID] = true
		
		// Find similar tickets within the same time window (1 hour)
		timeWindow := ticket.CreatedAt.Add(-1 * time.Hour)
		
		for j := i + 1; j < len(tickets); j++ {
			otherTicket := tickets[j]
			if processed[otherTicket.ID] {
				continue
			}
			
			// Check if within time window
			if otherTicket.CreatedAt.Before(timeWindow) {
				continue
			}
			
			// Simple similarity check - same type and similar title keywords
			if s.areTicketsSimilarSimpleSimilar(ticket, otherTicket) {
				group = append(group, otherTicket.ID)
				processed[otherTicket.ID] = true
			}
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

// areTicketsSimilarSimpleSimilar is a simplified similarity check for SimilarTrouble structs
func (s *AccumulationService) areTicketsSimilarSimpleSimilar(ticket1, ticket2 SimilarTrouble) bool {
	// Same type is a strong indicator
	if ticket1.Type == ticket2.Type && ticket1.Type != "" {
		return true
	}
	
	// Simple keyword overlap in titles
	title1Keywords := s.extractKeywords(ticket1.Title)
	title2Keywords := s.extractKeywords(ticket2.Title)
	
	commonKeywords := 0
	for _, kw1 := range title1Keywords {
		for _, kw2 := range title2Keywords {
			if kw1 == kw2 {
				commonKeywords++
				break
			}
		}
	}
	
	// If more than 1 common keyword, consider similar
	return commonKeywords >= 1
}
