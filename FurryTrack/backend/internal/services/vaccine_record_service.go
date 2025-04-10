package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	"fmt"
	"time"
	"log"

	"github.com/google/uuid"
)

type VaccineRecordService struct {
	recordRepo  *repositories.VaccineRecordRepository
	petRepo     *repositories.PetRepository
	vaccineRepo *repositories.VaccineRepository
}


func NewVaccineRecordService(
    recordRepo *repositories.VaccineRecordRepository,
    petRepo *repositories.PetRepository,
    vaccineRepo *repositories.VaccineRepository,
) *VaccineRecordService {
    return &VaccineRecordService{
        recordRepo:  recordRepo,
        petRepo:     petRepo,    
        vaccineRepo: vaccineRepo, 
    }
}

func (s *VaccineRecordService) AddVaccineRecord(
	userID uuid.UUID,
	petID uuid.UUID,
	vaccineName string,
	date time.Time,
	clinic string,
	userRole models.Role,
) (*models.VaccineRecord, error) {
	// Проверяем права доступа
	log.Printf("AddVaccineRecord - petRepo: %v", s.petRepo != nil)
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}
	if pet.OwnerID != userID && userRole != models.RoleVet{
		return nil, fmt.Errorf("access denied")
	}

	vaccine, err := s.vaccineRepo.FindByName(vaccineName)
	if err != nil {
		return nil, fmt.Errorf("vaccine not found: %w", err)
	}

	// Создаем запись с автоматическим расчетом следующей даты прививки
	record := &models.VaccineRecord{
		UserID:    uuid.New(),
		PetID:     petID,
		VaccineID: vaccine.ID,
		VaccineName: vaccineName,
		Date:      date,
		Clinic:    clinic,
		NextDate:  date.AddDate(0, 0, vaccine.DurationDays),
	}

	if err := s.recordRepo.AddRecord(record); err != nil {
		return nil, fmt.Errorf("failed to save record: %w", err)
	}

	return record, nil
}

// GetPetVaccinationHistory — получить историю прививок питомца
func (s *VaccineRecordService) GetPetVaccinationHistory(petID uuid.UUID) ([]models.VaccineRecord, error) {
	return s.recordRepo.GetRecordsByPet(petID)
}
