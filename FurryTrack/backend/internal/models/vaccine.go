package models

import (
	// "gorm.io/gorm"
	"github.com/google/uuid"
)

// Vaccine — справочник возможных вакцин
type Vaccine struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`          
	Description  string    `json:"description"`   
	DurationDays int       `json:"duration_days"`
}
