package repositories

import (
	"FurryTrack/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaccineRecordRepository struct {
	db *gorm.DB
}

func NewVaccineRecordRepository(db *gorm.DB) *VaccineRecordRepository {
	return &VaccineRecordRepository{db: db}
}

// AddRecord - добавление записи
func (r *VaccineRecordRepository) AddRecord(record *models.VaccineRecord) error {
	return r.db.Create(record).Error
}

// GetRecordsByPet - получения списка вакцин питомца
func (r *VaccineRecordRepository) GetRecordsByPet(petID uuid.UUID) ([]models.VaccineRecord, error) {
	var records []models.VaccineRecord
	err := r.db.
		Where("pet_id = ?", petID).
		Order("date DESC").
		Find(&records).Error
	return records, err
}

// GetByID - получение записи о вакцине по ID
func (r *VaccineRecordRepository) GetByID(id uuid.UUID) (*models.VaccineRecord, error) {
	var record models.VaccineRecord
	err := r.db.First(&record, "id = ?", id).Error
	return &record, err
}

// GetByID - получение записи о вакцине по ID пользователя
func (r *VaccineRecordRepository) GetByUserID(userID uuid.UUID) ([]models.VaccineRecord, error) {
	var records []models.VaccineRecord
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&records).Error
	return records, err
}
