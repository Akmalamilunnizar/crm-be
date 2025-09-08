package ticketapi

import (
	"log"
	"skripsi-be/internal/ml"
	"skripsi-be/internal/models/entities"
)

// MLService handles machine learning operations for trouble tickets
type MLService struct {
	classifier *ml.SVMClassifier
}

// NewMLService creates a new ML service
func NewMLService() *MLService {
	return &MLService{
		classifier: ml.NewSVMClassifier(),
	}
}

// ClassifyTroubleTicket classifies a trouble ticket based on its title
func (mls *MLService) ClassifyTroubleTicket(title string) (*ml.TroubleType, error) {
	log.Printf("Classifying trouble ticket: '%s'", title)

	result, err := mls.classifier.Predict(title)
	if err != nil {
		log.Printf("Error classifying trouble ticket: %v", err)
		return nil, err
	}

	log.Printf("Classification result: %s (confidence: %.2f)", result.Type, result.Confidence)
	return result, nil
}

// AutoClassifyTicket automatically classifies a ticket and updates its type
func (mls *MLService) AutoClassifyTicket(ticket *entities.TroubleTicket) (*entities.TroubleTicket, error) {
	if ticket.Title == "" {
		log.Println("Ticket title is empty, cannot classify")
		return ticket, nil
	}

	// Classify the ticket
	classification, err := mls.ClassifyTroubleTicket(ticket.Title)
	if err != nil {
		log.Printf("Failed to classify ticket: %v", err)
		return ticket, err
	}

	// Update ticket type if confidence is high enough
	if classification.Confidence > 0.6 { // Threshold for auto-assignment
		ticket.Type = &classification.Type
		log.Printf("Auto-assigned ticket type: %s (confidence: %.2f)", classification.Type, classification.Confidence)
	} else {
		log.Printf("Low confidence classification (%.2f), keeping manual assignment", classification.Confidence)
	}

	return ticket, nil
}

// GetClassificationStats returns statistics about the ML classifier
func (mls *MLService) GetClassificationStats() map[string]interface{} {
	return mls.classifier.GetClassificationStats()
}

// BatchClassifyTickets classifies multiple tickets at once
func (mls *MLService) BatchClassifyTickets(tickets []entities.TroubleTicket) ([]entities.TroubleTicket, error) {
	log.Printf("Batch classifying %d tickets", len(tickets))

	results := make([]entities.TroubleTicket, len(tickets))

	for i, ticket := range tickets {
		classifiedTicket, err := mls.AutoClassifyTicket(&ticket)
		if err != nil {
			log.Printf("Error classifying ticket %d: %v", i, err)
			results[i] = ticket // Keep original if classification fails
		} else {
			results[i] = *classifiedTicket
		}
	}

	return results, nil
}

// GetSuggestedTypes returns suggested trouble types based on title
func (mls *MLService) GetSuggestedTypes(title string) ([]ml.TroubleType, error) {
	if title == "" {
		return []ml.TroubleType{}, nil
	}

	// Get classification result
	result, err := mls.ClassifyTroubleTicket(title)
	if err != nil {
		return nil, err
	}

	// Return as single suggestion for now
	// In advanced implementation, could return multiple suggestions with confidence scores
	return []ml.TroubleType{*result}, nil
}
