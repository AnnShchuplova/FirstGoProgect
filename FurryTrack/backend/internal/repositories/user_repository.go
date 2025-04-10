package repositories

import (
	"FurryTrack/internal/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	//"golang.org/x/crypto/bcrypt"
	"strings"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// CreateUser - создаем нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return errors.New("user already exists")
		}
		return result.Error
	}
	return nil
}

// Находим пользователя по email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByID - находим пользователя по ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// ExistsByEmail - проверяет существование пользователя по email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil 
		}
		return false, fmt.Errorf("database error: %v", err) 
	}
	return true, nil 
}

// UpdateUser - обновление информации о пользователе
func (r *UserRepository) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
    result := r.db.Model(&models.User{}).
        Where("id = ?", userID).
        Updates(updates)

    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    
    return nil
}