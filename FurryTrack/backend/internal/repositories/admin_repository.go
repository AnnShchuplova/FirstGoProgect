package repositories

import (
	"FurryTrack/internal/models"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *AdminRepository) BanUser(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"banned_at": now,
			"role":      models.RoleUser,
		}).Error
}

func (r *AdminRepository) UnbanUser(userID uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("banned_at", nil).Error
}

func (r *AdminRepository) DeleteUser(userID uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", userID).Error
}

func (r *AdminRepository) LogAction(action *models.AdminAction) error {
	return r.db.Create(action).Error
}

func (r *AdminRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", userID).Error
	return &user, err
}