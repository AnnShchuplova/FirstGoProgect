package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"FurryTrack/internal/models"
	"FurryTrack/internal/services"
)

type EventController struct {
	service *services.EventService
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{service: service}
}

// CreateEvent создает новое событие
func (c *EventController) CreateEvent(ctx *gin.Context) {
	var input struct {
		PetID       uuid.UUID       `json:"pet_id" binding:"required"`
		Type        models.EventType `json:"type" binding:"required"`
		Title       string          `json:"title" binding:"required"`
		Description string          `json:"description"`
		Date        time.Time       `json:"date" binding:"required"`
		Location    string          `json:"location"`
		Cost        float64         `json:"cost"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)

	event := &models.Event{
		PetID:       input.PetID,
		UserID:      userID,
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Location:    input.Location,
		Cost:        input.Cost,
	}

	if err := c.service.CreateEvent(event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// GetPetEvents возвращает события для питомца
func (c *EventController) GetPetEvents(ctx *gin.Context) {
	petID, err := uuid.Parse(ctx.Param("pet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pet ID"})
		return
	}

	events, err := c.service.GetEventsByPetID(petID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": events,
		"meta": gin.H{
			"count": len(events),
		},
	})
}

// Обновление события
func (c *EventController) UpdateEvent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	var updateData models.Event
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	existingEvent, err := c.service.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	existingEvent.Title = updateData.Title
	existingEvent.Description = updateData.Description
	existingEvent.Date = updateData.Date
	existingEvent.Location = updateData.Location
	existingEvent.Cost = updateData.Cost
	existingEvent.Type = updateData.Type

	if err := c.service.UpdateEvent(existingEvent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	ctx.JSON(http.StatusOK, existingEvent)
}

// Удаление события
func (c *EventController) DeleteEvent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	// Проверяем существование события
	if _, err := c.service.GetEventByID(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := c.service.DeleteEvent(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	ctx.Status(http.StatusNoContent)
}