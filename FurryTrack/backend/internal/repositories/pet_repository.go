package repositories

import (
	"FurryTrack/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	//"errors"
	//"fmt"
	//"log"
)

type PetRepository struct {
	*BaseRepository
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

//func (r *PetRepository) checkDB() error {
//    if r.db == nil {
//        return fmt.Errorf("database connection is not initialized")
//    }
//    return nil
//}

// Create - создаем нового питомца
func (r *PetRepository) Create(pet *models.Pet) error {
	return r.db.Create(pet).Error
}

// FindByOwner - находим всех питомцев владельца
func (r *PetRepository) FindByOwner(ownerID uuid.UUID) ([]models.Pet, error) {
	var pets []models.Pet
	err := r.db.Where("owner_id = ?", ownerID).Find(&pets).Error
	return pets, err
}

// FindByID - находим питомца по ID
func (r *PetRepository) FindByID(id uuid.UUID) (*models.Pet, error) {
    
    var pet models.Pet
    err := r.db.Where("id = ?", id).First(&pet).Error
    if err != nil {

        return nil, err
    }
    
    return &pet, nil
}

func (r *PetRepository) Update(pet *models.Pet) error {
    return r.db.Save(pet).Error
}

func (r *PetRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Pet{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

// UpdatePhoto - обновление фотографии питомца
func (r *PetRepository) UpdatePhoto(petID uuid.UUID, photoURL string) error {
	return r.db.Model(&models.Pet{}).
		Where("id = ?", petID).
		Update("photo_url", photoURL).
		Error
}