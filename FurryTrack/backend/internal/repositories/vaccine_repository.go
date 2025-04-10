package repositories

import (
	"FurryTrack/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaccineRepository struct {
	*BaseRepository
}

func NewVaccineRepository(db *gorm.DB) *VaccineRepository {
	return &VaccineRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetAll - получение всего справочника вакцин
func (r *VaccineRepository) GetAll() ([]models.Vaccine, error) {
	var vaccines []models.Vaccine
	result := r.db.Find(&vaccines)
	if result.Error != nil {
		return nil, result.Error
	}
	return vaccines, nil
}

// FindByName - поиск вакцины по названию 
func (r *VaccineRepository) FindByName(name string) (*models.Vaccine, error) {
	var vaccine models.Vaccine
	result := r.db.Where("name = ?", name).First(&vaccine)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vaccine, nil
}

// CreateVaccine - создает новую вакцину в справочнике
func (r *VaccineRepository) CreateVaccine(vaccine *models.Vaccine) error {
	if vaccine.ID == uuid.Nil {
		vaccine.ID = uuid.New()
	}

	result := r.db.Create(vaccine)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID - получение вакцины по ID
func (r *VaccineRepository) GetByID(id uuid.UUID) (*models.Vaccine, error) {
	var vaccine models.Vaccine
	result := r.db.First(&vaccine, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vaccine, nil
}

// Update - обновление данных вакцины
func (r *VaccineRepository) Update(vaccine *models.Vaccine) error {
	result := r.db.Save(vaccine)
	return result.Error
}

// Delete - удаление вакцины
func (r *VaccineRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Vaccine{}, "id = ?", id)
	return result.Error
}
