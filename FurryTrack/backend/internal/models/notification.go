package models

import (
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Notification - модель данных для уведомлений
type Notification struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
    Type      string         `json:"type"`
    Message   string         `json:"message"`
    IsRead    bool           `json:"is_read" gorm:"default:false"`
    CreatedAt time.Time      `json:"created_at"`
    ExtraData string         `json:"extra_data"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}