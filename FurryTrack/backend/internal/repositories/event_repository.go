package repositories

import (
	"FurryTrack/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create - создает новое событие
func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

// GetByID - получает событие по ID 
func (r *EventRepository) GetByID(id uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("Pet").Preload("User").First(&event, "id = ?", id).Error
	return &event, err
}

// GetByPetID - получает все события питомца, отсортированные по дате (сначала новые)
func (r *EventRepository) GetByPetID(petID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	err := r.db.Where("pet_id = ?", petID).Order("date DESC").Find(&events).Error
	return events, err
}

// Update - обновляет данные события
func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

// Delete - удаляет событие
func (r *EventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Event{}, "id = ?", id).Error
}