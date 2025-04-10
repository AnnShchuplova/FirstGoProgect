package services

import (
    "encoding/json"
    
    "FurryTrack/internal/models"
    "FurryTrack/internal/repositories"
    
    "github.com/google/uuid"
)

type NotificationService struct {
    repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
    return &NotificationService{repo: repo}
}

// Обращения к репозиторию

func (s *NotificationService) CreateNotification(userID uuid.UUID, notifType, message string, extra interface{}) error {
    extraData, _ := json.Marshal(extra)
    
    notification := &models.Notification{
        UserID:    userID,
        Type:      notifType,
        Message:   message,
        ExtraData: string(extraData),
    }
    
    return s.repo.Create(notification)
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
    return s.repo.GetByUserID(userID)
}

func (s *NotificationService) MarkAsRead(id, userID uuid.UUID) error {
    return s.repo.MarkAsRead(id, userID)
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
    return s.repo.GetUnreadCount(userID)
}