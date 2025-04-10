package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User - модель пользователя
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Username     string    `gorm:"unique;not null"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	Role         Role `gorm:"type:varchar(10);not null;default:'USER'"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
    Pets         []Pet          `gorm:"foreignKey:OwnerID"`
    IsAdmin      bool           `gorm:"default:false"`
}

// LoginRequest - структура для входа пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse - структура ответа при успешном входе
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Для админа
type AdminAction struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AdminID     uuid.UUID   `gorm:"type:uuid" json:"admin_id"`
	UserID      uuid.UUID   `gorm:"type:uuid" json:"user_id"`
	ActionType  string      `json:"action_type"` // "ban", "unban", "delete", etc.
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
}