package repositories

import (
    "FurryTrack/internal/models"   
    "github.com/google/uuid"
    "gorm.io/gorm"
)


type NotificationRepository struct {
    db *gorm.DB 
}


func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
    return &NotificationRepository{db: db}
}

// Create - создает новое уведомление в базе данных
func (r *NotificationRepository) Create(notification *models.Notification) error {
    return r.db.Create(notification).Error
}

// GetByUserID - выводит все уведомления 
func (r *NotificationRepository) GetByUserID(userID uuid.UUID) ([]models.Notification, error) {
    var notifications []models.Notification
    err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
    return notifications, err
}

// MarkAsRead - помечает уведомление как прочитанное
func (r *NotificationRepository) MarkAsRead(id, userID uuid.UUID) error {
    return r.db.Model(&models.Notification{}).
        Where("id = ? AND user_id = ?", id, userID).
        Update("is_read", true).Error
}

// GetUnreadCount - возвращает количество непрочитанных уведомлений для пользователя
func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = false", userID).
        Count(&count).Error
    return count, err
}