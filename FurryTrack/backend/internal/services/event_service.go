package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	"github.com/google/uuid"
)

type EventService struct {
	repo *repositories.EventRepository 
}

func NewEventService(repo *repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// Обращения к репозиторию

func (s *EventService) CreateEvent(event *models.Event) error {
	return s.repo.Create(event)
}

func (s *EventService) GetEventByID(id uuid.UUID) (*models.Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) GetEventsByPetID(petID uuid.UUID) ([]models.Event, error) {
	return s.repo.GetByPetID(petID)
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	return s.repo.Update(event)
}

func (s *EventService) DeleteEvent(id uuid.UUID) error {
	return s.repo.Delete(id)
}