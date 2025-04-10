package models

import (
	"time"

	"github.com/google/uuid"
)

// VaccineRecord — отметка о проведённой вакцинации
type VaccineRecord struct {
	VaccineID   uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	PetID       uuid.UUID `json:"pet_id"`       
	VaccineName string    `json:"vaccine_name"` 
	Date        time.Time `json:"date"`         
	Clinic      string    `json:"clinic"`       
	NextDate    time.Time `json:"next_date"`    
}
