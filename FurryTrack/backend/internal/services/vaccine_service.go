package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	//"gorm.io/gorm"
	"github.com/google/uuid"
)

type VaccineService struct {
	vaccineRepo *repositories.VaccineRepository
}

func NewVaccineService(vaccineRepo *repositories.VaccineRepository) *VaccineService {
	return &VaccineService{
		vaccineRepo: vaccineRepo,
	}
}

// AddVaccine добавляет новую вакцину в справочник
func (s *VaccineService) AddVaccine(vaccine *models.Vaccine) error {
	return s.vaccineRepo.CreateVaccine(vaccine)
}

// GetAllVaccines возвращает список всех вакцин из справочника
func (s *VaccineService) GetAllVaccines() ([]models.Vaccine, error) {
	return s.vaccineRepo.GetAll()
}

// GetVaccineByID возвращает вакцину по её ID
func (s *VaccineService) GetVaccineByID(id uuid.UUID) (*models.Vaccine, error) {
	return s.vaccineRepo.GetByID(id)
}

// UpdateVaccine обновляет данные вакцины
func (s *VaccineService) UpdateVaccine(vaccine *models.Vaccine) error {
	return s.vaccineRepo.Update(vaccine)
}

// DeleteVaccine удаляет вакцину из справочника
func (s *VaccineService) DeleteVaccine(id uuid.UUID) error {
	return s.vaccineRepo.Delete(id)
}

// FindVaccineByName ищет вакцину по названию 
func (s *VaccineService) FindVaccineByName(name string) (*models.Vaccine, error) {
	return s.vaccineRepo.FindByName(name)
}
