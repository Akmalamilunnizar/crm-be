package handlers

import (
	"strconv"
	"time"

	"skripsi-be/internal/helpers"
	"skripsi-be/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AccumulationHandler struct {
	accumulationService *services.AccumulationService
}

func NewAccumulationHandler(accumulationService *services.AccumulationService) *AccumulationHandler {
	return &AccumulationHandler{
		accumulationService: accumulationService,
	}
}

// GetSimilarTroubles finds similar troubles for a given ticket
func (h *AccumulationHandler) GetSimilarTroubles(c *fiber.Ctx) error {
	ticketIDStr := c.Params("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 64)
	if err != nil {
		return helpers.ResponseUtils(c, 400, false, "Invalid ticket ID", nil)
	}
	
	timeWindowStr := c.Query("time_window", "24")
	timeWindow, err := strconv.Atoi(timeWindowStr)
	if err != nil {
		timeWindow = 24
	}
	
	similarTroubles, err := h.accumulationService.DetectSimilarTroubles(ticketID, timeWindow)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	
	return helpers.ResponseUtils(c, 200, true, "Similar troubles found", fiber.Map{
		"data": similarTroubles,
		"count": len(similarTroubles),
	})
}

// UpdateAccumulation updates the accumulation count for a group of tickets
func (h *AccumulationHandler) UpdateAccumulation(c *fiber.Ctx) error {
	var request struct {
		TicketIDs    []uint64 `json:"ticket_ids"`
		Accumulation int       `json:"accumulation"`
	}
	
	if err := c.BodyParser(&request); err != nil {
		return helpers.ResponseUtils(c, 400, false, err.Error(), nil)
	}
	
	if len(request.TicketIDs) == 0 {
		return helpers.ResponseUtils(c, 400, false, "ticket_ids is required", nil)
	}
	
	if request.Accumulation < 1 {
		return helpers.ResponseUtils(c, 400, false, "accumulation must be at least 1", nil)
	}
	
	err := h.accumulationService.UpdateAccumulation(request.TicketIDs, request.Accumulation)
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	
	return helpers.ResponseUtils(c, 200, true, "Accumulation updated successfully", fiber.Map{
		"ticket_ids": request.TicketIDs,
		"accumulation": request.Accumulation,
	})
}

// AutoDetectAndGroup automatically detects and groups similar troubles
func (h *AccumulationHandler) AutoDetectAndGroup(c *fiber.Ctx) error {
	err := h.accumulationService.AutoDetectAndGroup()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	
	return helpers.ResponseUtils(c, 200, true, "Automatic grouping completed", fiber.Map{
		"timestamp": time.Now(),
	})
}

// GetAccumulationStats returns statistics about ticket accumulation
func (h *AccumulationHandler) GetAccumulationStats(c *fiber.Ctx) error {
	stats, err := h.accumulationService.GetAccumulationStats()
	if err != nil {
		return helpers.ResponseUtils(c, 500, false, err.Error(), nil)
	}
	
	return helpers.ResponseUtils(c, 200, true, "Accumulation statistics", stats)
}

// GetHighAccumulationTickets returns tickets with high accumulation (multiple customers affected)
func (h *AccumulationHandler) GetHighAccumulationTickets(c *fiber.Ctx) error {
	minAccumulationStr := c.Query("min_accumulation", "2")
	minAccumulation, err := strconv.Atoi(minAccumulationStr)
	if err != nil {
		minAccumulation = 2
	}
	
	// This would need to be implemented in the service
	// For now, return a placeholder response
	return helpers.ResponseUtils(c, 200, true, "High accumulation tickets endpoint - to be implemented", fiber.Map{
		"data": []interface{}{},
		"min_accumulation": minAccumulation,
	})
}
