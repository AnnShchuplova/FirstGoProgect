package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type EventType string

const (
	EventVetVisit    EventType = "vet_visit"
	EventVaccination EventType = "vaccination"
	EventMedication  EventType = "medication"
	EventOther       EventType = "other"
)

// Event - модель данных для события
type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PetID       uuid.UUID `gorm:"type:uuid;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Type        EventType `gorm:"type:varchar(50);not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"not null"`
	Location    string    `gorm:"type:varchar(200)"`
	Cost        float64   `gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Pet  Pet  `gorm:"foreignKey:PetID"`
	User User `gorm:"foreignKey:UserID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}