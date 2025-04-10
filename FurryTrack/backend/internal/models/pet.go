package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Pet - модель данных для питомца
type Pet struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"pet_id"`
    Name      string    `gorm:"type:text;not null"`
    Type      string    `gorm:"type:text;not null"`
    Breed     string    `gorm:"type:text"`
    BirthDate time.Time `gorm:"type:date"` 
    OwnerID   uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
    CreatedAt time.Time `gorm:"type:timestamp without time zone;default:now()"`
    UpdatedAt time.Time `gorm:"type:timestamp without time zone"`
    DeletedAt gorm.DeletedAt `gorm:"type:timestamp without time zone;index"`
    PhotoURL  string         `gorm:"type:varchar(255)"`
}

// Модель данных для ветинарного визита, заменена общей моделью событий Event
type VetVisit struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PetID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"pet_id"`
	Date        string
	Description string
}
